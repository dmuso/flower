api: nix-shell --run "cd api && air --build.cmd 'go build -o ./tmp/server ./cmd/server' --build.bin ./tmp/server --build.exclude_dir tmp --build.poll=500ms"
tailwind: nix-shell --run "cd frontend && bun run tailwind:watch"
frontend: nix-shell --run "cd frontend && FRONTEND_API_URL=${FRONTEND_API_URL} FRONTEND_PORT=${FRONTEND_PORT} bun --hot public/index.html --spa --host=0.0.0.0 --port=${FRONTEND_PORT}"
