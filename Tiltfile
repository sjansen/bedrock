load('ext://restart_process', 'docker_build_with_restart')


local_resource(
    'bedrock-compile',
    'CGO_ENABLED=0 GOOS=linux go build -o dist/bedrock',
    deps=['./main.go', './internal/'],
    ignore=[
        '**/*_test.go',
        '**/*.templ',
        '**/*.templ*', # created by templ fmt
    ],
    resource_deps=['bedrock-generate'],
)

local_resource(
    'bedrock-generate',
    'go generate',
    deps=['./internal/templates/'],
    ignore=['**/*_templ.go'],
)

local_resource(
    'tailwindcss',
    'tailwindcss -i public/default.css -o dist/public/default.css --minify',
    deps=['./public/'],
    ignore=['**', '!public/**/*.css']
)

docker_build_with_restart(
    'bedrock',
    '.', 
    dockerfile='localdev/Dockerfile',
    entrypoint=['/app/bedrock'],
    live_update=[
        sync('./dist/bedrock', '/app/bedrock'),
        sync('./dist/public', '/app/public'),
    ],
    only=['./dist'],
)

k8s_yaml('localdev/app.yaml')
k8s_resource(
    'bedrock',
    port_forwards=8080,
    resource_deps=[
        'bedrock-compile',
        'bedrock-generate',
        'tailwindcss',
    ],
)
