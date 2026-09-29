import React from "react";
import { CategoryTag, SeverityBadge, StatusBadge } from "./Badges";
import { formatDate } from "../utils";

export function ReportCard({ report }) {
  return (
    <article className="card report-card">
      <header>
        <h2>
          <a href={`#/reports/${report.id}`}>{report.title}</a>
        </h2>
        <StatusBadge status={report.status} />
      </header>
      <div className="row">
        <CategoryTag category={report.category} />
        <SeverityBadge severity={report.severity} />
        <span className="meta">{formatDate(report.created_at)}</span>
      </div>
      {report.description && <p>{report.description}</p>}
      {report.photo_url && <img className="photo-preview" src={report.photo_url} alt="" />}
    </article>
  );
}
