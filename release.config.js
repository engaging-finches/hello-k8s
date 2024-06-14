const config={
    branches: ['main'], //only release on main branch
    plugins: [
        '@semantic-release/commit-analyzer', //added by default
        {
            "preset": "conventionalcommits"
        },
        '@semantic-release/release-notes-generator', //added by default
        '@semantic-release/npm', //added by default
        '@semantic-release/github', //added by default
        "@semantic-release/changelog",
        ["@semantic-release/git", {
            "assets": ["version.txt", "CHANGELOG.md"], //commits these files back into the repo
            "message": "chore(release): ${nextRelease.version} [skip ci]\n\n${nextRelease.notes}" //skip commits with [skip ci]
        }],
    ]

}

module.exports=config;