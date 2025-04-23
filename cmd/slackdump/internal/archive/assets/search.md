# Search Command

The `search` command allows you to search and export messages or files from a
Slack workspace based on specified query terms. This command supports searching
for messages, files, or both, and outputs the results in a directory.

### Subcommands
- **`slackdump search messages <query>`**: Searches and records messages
  matching the given query.
- **`slackdump search files <query>`**: Searches and records files matching the given
  query.
- **`slackdump search all <query>`**: Searches and records both messages and files
  matching the query.

### Flags
- **`-no-channel-users`**: Skips retrieving user data for channels, making the
  process approximately 2.5x faster.

### Requirements
- Authentication is required for all search operations.
- An output directory must be specified (see configuration details).

## Usage Examples

### Search Messages

```bash
slackdump search messages "meeting notes"
```

### Search Files

```bash
slackdump search files "report"
```

### Search Messages and Files

```bash
slackdump search all "project updates"
```

### Faster Searches
To speed up searches, add the `-no-channel-users` flag:

```bash
slackdump search messages -no-channel-users "status update"
```


## Output
The search command outputs results to the database in the output directory.
See `slackdump help archive` for details on the database structure and
contents.
