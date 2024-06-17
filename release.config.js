const config = {
  branches: ['main', 'edit-semantic'], //only release on main branch
  dryRun: false, //set to false to actually release
  debug: false, //set to false to turn off debugging
  plugins: [
    '@semantic-release/commit-analyzer', //added by default
    {
      "preset": "conventionalcommits"
    },
    '@semantic-release/release-notes-generator', //added by default
    ['@semantic-release/github', {
      assets: ['ghrunner/dist/install.yaml'], // Include only 'dist/install.yaml' in the release
    }],
  ],
  repositoryUrl: "https://github.com/engaging-finches/hello-k8s",
}

module.exports = config;