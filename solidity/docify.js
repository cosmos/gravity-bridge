const NODE_DIR     = "./node_modules";
const INPUT_DIR    = "./contracts";
const CONFIG_DIR   = "./docgen";
const OUTPUT_DIR   = "./docgen/docs";
const README_FILE  = "./docgen/README.md";
const SUMMARY_FILE = "./docgen/SUMMARY.md";
const EXCLUDE_FILE = "./docgen/exclude.txt";

const fs        = require("fs");
const path      = require("path");
const spawnSync = require("child_process").spawnSync;

const excludeList  = lines(EXCLUDE_FILE).map(line => INPUT_DIR + "/" + line);
const relativePath = path.relative(path.dirname(SUMMARY_FILE), OUTPUT_DIR);

function lines(pathName) {
    return fs.readFileSync(pathName, {encoding: "utf8"}).split("\r").join("").split("\n");
}

function scan(pathName, indentation) {
    if (!excludeList.includes(pathName)) {
        if (fs.lstatSync(pathName).isDirectory()) {
            fs.appendFileSync(SUMMARY_FILE, indentation + "* " + path.basename(pathName) + "\n");
            for (const fileName of fs.readdirSync(pathName))
                scan(pathName + "/" + fileName, indentation + "  ");
        }
        else if (pathName.endsWith(".sol")) {
            const text = path.basename(pathName).slice(0, -4);
            const link = pathName.slice(INPUT_DIR.length, -4);
            fs.appendFileSync(SUMMARY_FILE, indentation + "* [" + text + "](" + relativePath + link + ".md)\n");
        }
    }
}

function fix(pathName) {
    if (fs.lstatSync(pathName).isDirectory()) {
        for (const fileName of fs.readdirSync(pathName))
            fix(pathName + "/" + fileName);
    }
    else if (pathName.endsWith(".md")) {
        fs.writeFileSync(pathName, lines(pathName).filter(line => line.trim().length > 0).join("\n\n") + "\n");
    }
}

fs.writeFileSync (SUMMARY_FILE, "# Summary\n");
fs.writeFileSync (".gitbook.yaml", "root: ./\n");
fs.appendFileSync(".gitbook.yaml", "structure:\n");
fs.appendFileSync(".gitbook.yaml", "  readme: " + README_FILE + "\n");
fs.appendFileSync(".gitbook.yaml", "  summary: " + SUMMARY_FILE + "\n");

scan(INPUT_DIR, "");

const args = [
    NODE_DIR + "/solidity-docgen/dist/cli.js",
    "--input="         + INPUT_DIR,
    "--output="        + OUTPUT_DIR,
    "--templates="     + CONFIG_DIR,
    "--solc-module="   + NODE_DIR + "/solc",
    "--solc-settings=" + JSON.stringify({optimizer: {enabled: true, runs: 200}})
];

const result = spawnSync("node", args, {stdio: ["inherit", "inherit", "pipe"]});
if (result.stderr.length > 0)
    throw new Error(result.stderr);

fix(OUTPUT_DIR);