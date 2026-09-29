import { titleCase, classNames } from "../utils";

export function SeverityBadge({ severity }) {
  return <span className={classNames("badge", `sev-${severity}`)}>Severity {severity}</span>;
}

export function CategoryTag({ category }) {
  return <span className="tag">{titleCase(category)}</span>;
}

export function StatusBadge({ status }) {
  return <span className={classNames("badge", `status-${status}`)}>{titleCase(status)}</span>;
}
