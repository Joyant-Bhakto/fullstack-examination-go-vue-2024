// @ts-check
import withNuxt from "./eslint.config.mjs";

export default withNuxt({
  files: ["ui/**"],
  rules: {
    "@typescript-eslint/explicit-function-return-type": "error",
  },
});
