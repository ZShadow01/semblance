const path = require("path");
const fs = require("fs");

/**
 * Iterate through all files and subdirectories of a given directory
 * @param {string} dir - The root directory
 * @param {(path: string, lstat: fs.Stats) => void} callback - The callback function
 */
function walkDir(dir, callback) {
    const lstat = fs.lstatSync(dir);
    const directories = [[dir, lstat]];

    if (lstat.isFile()) {
        callback(dir, lstat);
        return;
    }

    while (directories.length > 0) {
        const currentDir = directories.pop();

        callback(...currentDir);

        const files = fs.readdirSync(currentDir[0]);
        for (const file of files) {
            const filePath = path.join(currentDir[0], file);

            const lstat = fs.lstatSync(filePath);
            if (lstat.isDirectory()) {
                directories.push([filePath, lstat]);
            } else {
                callback(filePath, lstat);
            }
        }
    }
}

module.exports = { walkDir };
