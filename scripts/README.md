# Building the Application

The build out of Sidebar comes through a bash build script in this directory.
To start the build, run the following command from the project root.

```bash
# Script help output:
# Must provide build mode as argument
# Options are:
#   build              - build backend + vite + copy binary, no packaging
#   dev                - backend + vite dev server + electron
#   release-build      - full backend + vite build + electron-builder

./scripts/build.sh dev
```

# What to do when

> TODO: Add more information here...

Building to `dev` will be for frontend work and testing.
Building to `build` will be for backend work and testing.
Building to `release-build` will be for releasing the binaries for user install.
