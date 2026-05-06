# General

This strapi plugin use [yalc](https://www.npmjs.com/package/yalc) to symlink the package to the strapi project, so you can test it without publish it to npm.

This project use go-task to run utility commands. Installation of this tool is required.

For macOS, you can install it using homebrew:

```bash
brew install go-task
```

View available commands

```bash
task --list
```

# Development

First you need to install yalc globally

```bash
task install-yalc
```

Then you can publish the plugin to local yalc registry

```bash
task publish-plugin
```

After that, you can install the plugin in your strapi project

```bash
cd your-strapi-project
pnpm dlx yalc add --link cheap-strapi-plugin && pnpm install
```


# Resource references
- [yalc](https://www.npmjs.com/package/yalc)
- [go-task](https://taskfile.dev/)
- [Strapi plugin development documentation](https://docs.strapi.io/cms/plugins-development/create-a-plugin)