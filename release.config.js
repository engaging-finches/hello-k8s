const config={
    branches: ['main'], //only release on main branch
    dryRun: true, //set to false to actually release
    debug: false, //set to false to turn off debugging
    plugins: [
        '@semantic-release/commit-analyzer', //added by default
        {
            "preset": "conventionalcommits"
        },
        '@semantic-release/release-notes-generator', //added by default
        '@semantic-release/github', //added by default
    ]

}

module.exports=config;