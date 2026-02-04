module.exports = {
  branches: ['main', 'next', 'next-major', { name: 'beta', prerelease: true }, { name: 'alpha', prerelease: true }],
  dryRun: false,
  plugins: [
    "@semantic-release/commit-analyzer",
    "@semantic-release/release-notes-generator",
    [
      "@semantic-release/changelog",
      {
        changelogFile: "CHANGELOG.md",
      },
    ],
    "@semantic-release/github",
    [
      "@semantic-release/exec",
      {
        prepareCmd: "bash scripts/bump_version.sh ${nextRelease.version}",
      },
    ],
    [
      "@semantic-release/git",
      {
        assets: [
          "CHANGELOG.md",
          "VERSION",
          "VERIFICATION.txt",
          "Makefile",
          "debian/rules"
        ],
        message:
          "chore(release): ${nextRelease.version} [skip ci]\n\n${nextRelease.notes}",
      },
    ],
  ],
};
