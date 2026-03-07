import CreateExperimentError from "./CreateExperimentError";
import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";

describe("CreateExperimentError", () => {
  it("should render nothing when message is null", () => {
    const { container } = render(<CreateExperimentError message={null} />);
    expect(container.firstChild).toBeNull();
  });

  it("should render nothing when message is empty string", () => {
    const { container } = render(<CreateExperimentError message="" />);
    expect(container.firstChild).toBeNull();
  });

  it("should render error message when message is provided", () => {
    const errorMessage = "Failed to create experiment";
    render(<CreateExperimentError message={errorMessage} />);

    expect(screen.getByText("Something Went Wrong:")).toBeInTheDocument();
    expect(screen.getByText(errorMessage)).toBeInTheDocument();
  });
});
