# Pushing client binaries

1. Prepare and push changes to main
2. Check [releases](https://github.com/Johnkhk/cli_chat_app/releases) or [tags](https://github.com/Johnkhk/cli_chat_app/tags)
3. run `git tag -a v1.0.0 -m "Release v1.0.0: Initial release of client binaries"`
4. Push to tag `git push origin v1.0.0` -- the gh action should deploy it


```bash
# Delete the local tag
git tag -d v1.0.0

# Delete the remote tag
git push origin :refs/tags/v1.0.0

# Recreate the tag with the changes
git tag -a v1.0.0 -m "Release v1.0.0: Updated with CGO fix"

# Push the new tag
git push origin v1.0.0
```
