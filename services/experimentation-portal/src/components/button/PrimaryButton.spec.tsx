import PrimaryButton from "./PrimaryButton";
import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";

describe("PrimaryButton", () => {
  it("should render children content", () => {
    render(<PrimaryButton>Click me</PrimaryButton>);
    expect(screen.getByText("Click me")).toBeInTheDocument();
  });

  it("should call onClick handler when clicked", async () => {
    const handleClick = vi.fn();
    const user = userEvent.setup();

    render(<PrimaryButton onClick={handleClick}>Click me</PrimaryButton>);

    await user.click(screen.getByText("Click me"));

    expect(handleClick).toHaveBeenCalledTimes(1);
  });

  it("should merge custom className with default classes", () => {
    render(<PrimaryButton className="custom-class">Button</PrimaryButton>);

    const button = screen.getByText("Button");
    expect(button.className).toContain("custom-class");
    expect(button.className).toContain("hover:bg-blue-600");
    expect(button.className).toContain("rounded-xl");
  });

  it("should pass through additional button attributes", () => {
    render(
      <PrimaryButton type="submit" aria-label="Submit form">
        Submit
      </PrimaryButton>,
    );

    const button = screen.getByText("Submit");
    expect(button).toHaveAttribute("type", "submit");
    expect(button).toHaveAttribute("aria-label", "Submit form");
  });
});
