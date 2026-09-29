import React from "react";
import { classNames } from "../utils";
import { formatDate } from "../utils";

export function CommentThread({ comments = [] }) {
  if (!comments.length) return <p className="muted">No comments yet.</p>;
  return comments.map((comment) => (
    <div
      key={comment.id}
      className={classNames("comment", comment.is_authority_comment && "authority")}
    >
      <p className="meta">
        <strong>{comment.user?.name || "User"}</strong>
        {comment.is_authority_comment ? " · Authority" : ""} · {formatDate(comment.created_at)}
      </p>
      <p>{comment.body}</p>
    </div>
  ));
}
