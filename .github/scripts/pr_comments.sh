#!/usr/bin/env bash
set -euo pipefail

if [[ -z "${GITHUB_USER:-}" ]]; then
    echo "GITHUB_USER env variable is not set"
    exit 1
fi

if [[ -z "${GH_TOKEN:-}" ]]; then
    echo "GH_TOKEN env variable is not set"
    exit 1
fi

if [[ -z "${GITHUB_REPOSITORY:-}" ]]; then
    echo "GITHUB_REPOSITORY env variable is not set"
    exit 1
fi

if [[ -z "${PR_NUMBER:-}" ]]; then
    echo "PR_NUMBER env variable is not set"
    exit 1
fi

if [[ -z "${MAGIC_COMMENT_HINT:-}" ]]; then
    echo "MAGIC_COMMENT_HINT env variable is not set"
    exit 1
fi

cleanup_comment() {
    magic_comment_id=$(get_magic_comment_id)
    if [[ -n "$magic_comment_id" ]]; then
        echo "Cleaning up comment with id ${magic_comment_id}" >&2
        write_comment "Comment cleaned up :heavy_check_mark:"
        sleep 10
        echo "Deleting comment with id ${magic_comment_id}..." >&2
        gh api --method DELETE "repos/${GITHUB_REPOSITORY}/issues/comments/${magic_comment_id}"
    fi
}

write_comment() {
    local body
    body=$(printf '%s\n%s' "$MAGIC_COMMENT_HINT" "$1")
    magic_comment_id=$(get_magic_comment_id)

    if [[ -z "$magic_comment_id" ]]; then
        echo "Creating new comment..." >&2
        gh api --method POST "repos/${GITHUB_REPOSITORY}/issues/${PR_NUMBER}/comments" -f body="$body"
    else
        echo "Updating existing comment with id ${magic_comment_id}..." >&2
        gh api --method PATCH "repos/${GITHUB_REPOSITORY}/issues/comments/${magic_comment_id}" -f body="$body"
    fi
}

get_magic_comment_id() {
    echo "Checking for a comment with magic hint ${MAGIC_COMMENT_HINT}..." >&2
    gh api "repos/${GITHUB_REPOSITORY}/issues/${PR_NUMBER}/comments?per_page=100" \
        --jq ".[] | select(.body | startswith(\"${MAGIC_COMMENT_HINT}\")) | .id" | head -n 1 || true
}

# Argument parsing
comment_text=
while :; do
    case $1 in
        --cleanup-comment)
            cleanup_comment
            break
            ;;
        --write-comment=?*)
            comment_text=${1#*=}
            write_comment "$comment_text"
            break
            ;;
        -?*)
            printf 'ERROR: Unknown option: %s\n' "$1" >&2
            exit 1
            ;;
        *)
            printf 'ERROR: No option specified\n' >&2
            exit 1
    esac
    shift
done
