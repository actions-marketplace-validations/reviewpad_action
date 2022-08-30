![Reviewpad-Background-Logo-Shorter@1 5x](https://user-images.githubusercontent.com/38539/185982194-21bf7bb2-d2d2-40ed-8893-82a216d848a5.png)

# Reviewpad GitHub Action 
[![x-ray-badge](https://img.shields.io/badge/Time%20to%20Merge-Strong%20team-ee9b00?link=https://xray.reviewpad.com/analysis?repository=https%3A%2F%2Fgithub.com%2Freviewpad%2Faction&style=plastic.svg)](https://xray.reviewpad.com/analysis?repository=https%3A%2F%2Fgithub.com%2Freviewpad%2Faction) [![CIDeploy](https://github.com/reviewpad/action/actions/workflows/cideploy.yml/badge.svg)](https://github.com/reviewpad/action/actions/workflows/cideploy.yml)

🔥 **Latest Stable Version**: v3.x ([Faro](https://en.wikipedia.org/wiki/Faro,_Portugal) Edition)

🤔 For **questions**, check out the [discussions](https://github.com/reviewpad/reviewpad/discussions).

📃 For **documentation**, check out this document and the [official documentation](https://docs.reviewpad.com).

🙌 **Join our community on [discord](https://reviewpad.com/discord)!**

___


This action runs the docker image [reviewpad/action](https://hub.docker.com/repository/docker/reviewpad/action).

It reads and automates the pull request workflows specified in the `reviewpad.yml` configuration file.

These workflows can be used to automatically label, assign reviewers, comment, merge and close pull requests.

For example, the following `reviewpad.yml` file:

```yaml
api-version: reviewpad.com/v3.x

workflows:
  - name: label-pull-request-with-size
    if:
      - rule: $size() <= 50
        extra-actions:
          - $addLabel("small")
      - rule: $size() > 50 && $size() <= 150
        extra-actions:
          - $addLabel("medium")
      - rule: $size() > 150
        extra-actions:
          - $addLabel("large")
```

Specifies a workflow to automatically add a label based on the size of the pull request.

For more information on the release procedure, check the [RELEASE.md](./RELEASE.md) document.


## Inputs

- **event**: The GitHub event context that trigger the action. Uses default `${{ toJSON(github) }}`
- **file**: The local location of the Reviewpad configuration file. Uses default `./reviewpad.yml`. Ignored if `file_url` is set.
- **file_url** *(OPTIONAL)*: The remote location of the Reviewpad configuration file. If set, it will ignore the input `file`.
- **token**: Uses default `${{ github.token }}`

## Usage

### Basic

Add the following step to your GitHub Action job:

```yaml
- name: Run reviewpad action
  uses: reviewpad/action@v3.x
```

### :link: Remote configuration

You can run reviewpad action with a remote configuration by setting the input `file_url`:

```yaml
- name: Run reviewpad action
  uses: reviewpad/action@v3.x
  with:
    file_url: https://github.com/reviewpad/action/blob/main/templates/start.yml
```

### :key: GitHub token

By default this action uses the `github-actions[bot]` token.

As described in the [official GitHub documentation](https://docs.github.com/en/actions/security-guides/automatic-token-authentication#using-the-github_token-in-a-workflow):

> When you use the repository's GITHUB_TOKEN to perform tasks, events triggered by the GITHUB_TOKEN will not create a new workflow run.

If you want to use more advanced features such as the [merge](https://docs.reviewpad.com/guides/built-ins#merge) feature, we recommend that you explicitly provide a PAT (Personal Access Token) to run this action:

```yaml
- name: Run reviewpad action
  uses: reviewpad/action@v3.x
  with:
    token: ${{ secrets.GH_TOKEN }}
```

[Please follow this link to know more](https://docs.reviewpad.com/getting-started/installation-with-github-token).
