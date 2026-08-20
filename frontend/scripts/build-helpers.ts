export function rewriteIndexHTMLContents(original: string): string {
  const scriptPlaceholder = "../src/index.tsx";
  const builtEntry = "/index.js";
  const cssHref = "/index.css";

  let updated = original.includes(scriptPlaceholder)
    ? original.replace(scriptPlaceholder, builtEntry)
    : original;

  if (!updated.includes(`href="${cssHref}"`)) {
    const cssLinkLine = `        <link rel="stylesheet" href="${cssHref}" />`;
    if (updated.includes("</head>")) {
      updated = updated.replace("</head>", `${cssLinkLine}\n    </head>`);
    } else {
      updated = `${updated}\n${cssLinkLine}`;
    }
  }

  return updated;
}
