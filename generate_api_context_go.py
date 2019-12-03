#!/usr/bin/env python3
import argparse
import os
import os.path
import glob
import re
import mmap

IMPORTS_FMT = \
"""import (
    "ezxray"
    "context"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/request"
    "github.com/aws/aws-sdk-go/service/{service_package_name}"
)

"""

PAGING_FUNC_RE = \
        re.compile(", fn func\(\*[^ ]*, bool\) bool, ")

XRAY_NOCONTEXT_WRAPPER_FMT = \
"""{withcontext_entire_signature}
    f := (*{service_package_name}.{client_type}).{op}WithContext
    mergedContext := ezxray.WithXrayContext(ctx)
    return f((*{service_package_name}.{client_type})(c), mergedContext, input{other_param_name_list}, opts...)
}}

func (c *{client_type}) {op}({input_expr}{other_param_list}) {suffix}
    return c.{op}WithContext(context.Background(), input{other_param_name_list})
}}

"""

SERVICE_DEFINITIONS_FILE_FMT = \
"""package {service_package_name}

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/client"
    "github.com/aws/aws-sdk-go/service/{service_package_name}"
)

type {client_type} {service_package_name}.{client_type}

func New(p client.ConfigProvider, cfgs ...*aws.Config) *{client_type} {{ return (*{client_type})({service_package_name}.New(p, cfgs...)) }}
"""

EXPOSED_APIS_WITH_CONTEXT = re.compile(b"""(^func \(c \*)([^ ]*)(\) )([^ ]*)(WithContext\(ctx aws\.Context, )(input \*[^ ]*)(, .*)(opts \.\.\.request\.Option\))( .*{$)""", re.MULTILINE)


def emit(output_file, contents):
    output_file.write(contents)


def generate_withcontext_wrappers(service_package_name, api_file, output_file):
    emit(output_file, f'package {service_package_name}\n\n')
    emit(output_file, IMPORTS_FMT.format(**locals()))

    data = mmap.mmap(api_file.fileno(), 0, access=mmap.ACCESS_READ)
    for withcontext_decl in re.finditer(EXPOSED_APIS_WITH_CONTEXT, data):
        withcontext_entire_signature, client_type, op, input_expr, other_params, suffix = (c.decode('utf-8') for c in withcontext_decl.group(0, 2, 4, 6, 7, 9))
        other_param_name_list = ''
        other_param_list = ''
        if other_params != ', ':
            assert re.match(PAGING_FUNC_RE, other_params), f'other param list "{other_params}" did not match expected regex'
            other_param_name_list = ', fn'
            assert other_params.endswith(', ')
            other_param_list = other_params[:-2]

        emit(output_file, XRAY_NOCONTEXT_WRAPPER_FMT.format(**locals()))

    return client_type


def generate_service_definitions(service_package_name, client_type, output_file):
    emit(output_file, SERVICE_DEFINITIONS_FILE_FMT.format(**locals()))

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("sdk", help="path to the AWS SDK for Go", type=str)
    args = parser.parse_args()

    apis_glob = os.path.join(args.sdk, 'service/*/api.go')
    api_gos = glob.glob(apis_glob)
    if not api_gos:
        print('no service/*/api.go files were found in sdk')
        return -1

    for api_go_path in api_gos:
        print(f'generating wrapper override for {api_go_path}')
        service_package_name = os.path.basename(os.path.dirname(api_go_path))
        with open(api_go_path, 'r') as f:
            output_dir = os.path.join('github.com', 'aws', 'aws-sdk-go', 'service', service_package_name)
            os.makedirs(output_dir, exist_ok=True)
            output_api_path = os.path.join(output_dir, 'xray_api.go')
            with open(output_api_path, 'w') as o:
                client_type = generate_withcontext_wrappers(service_package_name, f, o)

            output_service_path = os.path.join(output_dir, 'xray_service.go')
            with open(output_service_path, 'w') as o:
                generate_service_definitions(service_package_name, client_type, o)


if __name__ == '__main__':
    exit(main())
