import React from "react";

export class ErrorBoundary extends React.Component {
  constructor(props) {
    super(props);
    this.state = { error: null };
  }

  static getDerivedStateFromError(error) {
    return { error };
  }

  componentDidCatch(error, info) {
    console.error("UI error:", error, info);
  }

  render() {
    if (this.state.error) {
      return (
        <section className="card stack">
          <h1>Something went wrong</h1>
          <p className="alert error">{String(this.state.error?.message || this.state.error)}</p>
          <button
            className="btn primary"
            type="button"
            onClick={() => {
              this.setState({ error: null });
              window.location.hash = "#/feed";
            }}
          >
            Back to reports
          </button>
        </section>
      );
    }
    return this.props.children;
  }
}
