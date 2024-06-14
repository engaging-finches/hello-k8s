const config={
    branches: ['main'], //only release on main branch
    plugins: [
        '@semantic-release/commit-analyzer',
        '@semantic-release/release-notes-generator',
        ["@semantic-release/git", {
            "assets": ["/ghrunner", "/user_guides", "/docker", "package-lock.json", "package.json", "README.md",],
            "message": "chore(release): ${nextRelease.version} [skip ci]\n\n${nextRelease.notes}"
        }],
        '@semantic-release/github'
    ]

}

module.exports=config;