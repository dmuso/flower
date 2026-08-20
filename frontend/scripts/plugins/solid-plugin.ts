import { SolidPlugin } from "@dschz/bun-plugin-solid";

const pluginOptions = {
  generate: "dom",
  hydratable: true,
  babel: {
    presets: [
      ["solid", { generate: "dom", hydratable: true }],
      ["@babel/preset-typescript", { allowDeclareFields: true, onlyRemoveTypeImports: true }],
    ],
  },
} as Parameters<typeof SolidPlugin>[0];

export const solidPlugin = SolidPlugin(pluginOptions);

export default solidPlugin;
