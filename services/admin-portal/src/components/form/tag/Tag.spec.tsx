import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import Tag from "./Tag";

describe("Tag", () => {
  it("should add a tag when comma is pressed", async () => {
    const user = userEvent.setup();

    render(<Tag name="tags" />);

    const input = screen.getByRole("textbox");
    await user.type(input, "newtag,");

    expect(screen.getByText("newtag")).toBeInTheDocument();
    expect(input).toHaveValue("");
  });

  it("should not add duplicate tags", async () => {
    const user = userEvent.setup();

    render(<Tag name="tags" />);

    const input = screen.getByRole("textbox");
    await user.type(input, "tag1,");
    await user.type(input, "tag1,");

    const tags = screen.getAllByText("tag1");
    expect(tags).toHaveLength(1);
  });

  it("should show an error for tags with spaces", async () => {
    const user = userEvent.setup();

    render(<Tag name="tags" />);

    const input = screen.getByRole("textbox");
    await user.type(input, "tag with spaces,");

    expect(screen.getByText(/cannot contain spaces/i)).toBeInTheDocument();
    expect(screen.queryByText("tag with spaces")).not.toBeInTheDocument();
  });

  it("should clear error when user continues typing", async () => {
    const user = userEvent.setup();

    render(<Tag name="tags" />);

    const input = screen.getByRole("textbox");
    await user.type(input, "tag with spaces,");

    expect(screen.getByText(/cannot contain spaces/i)).toBeInTheDocument();

    await user.type(input, "a");

    expect(
      screen.queryByText(/cannot contain spaces/i),
    ).not.toBeInTheDocument();
  });

  it("should remove a tag when the cross icon is clicked", async () => {
    const user = userEvent.setup();

    render(<Tag name="tags" />);

    const input = screen.getByRole("textbox");
    await user.type(input, "tag1,");
    await user.type(input, "tag2,");

    expect(screen.getByText("tag1")).toBeInTheDocument();
    expect(screen.getByText("tag2")).toBeInTheDocument();

    const removeButtons = screen.getAllByRole("button");
    await user.click(removeButtons[0]);

    expect(screen.queryByText("tag1")).not.toBeInTheDocument();
    expect(screen.getByText("tag2")).toBeInTheDocument();
  });

  it("should clear input after adding a tag", async () => {
    const user = userEvent.setup();

    render(<Tag name="tags" />);

    const input = screen.getByRole("textbox");
    await user.type(input, "newtag,");

    expect(input).toHaveValue("");
  });

  it("should not add empty tags", async () => {
    const user = userEvent.setup();

    render(<Tag name="tags" />);

    const input = screen.getByRole("textbox");
    await user.type(input, ",");

    const tagBubbles = screen.queryAllByRole("button");
    expect(tagBubbles).toHaveLength(0);
  });

  it("should display multiple tags", async () => {
    const user = userEvent.setup();

    render(<Tag name="tags" />);

    const input = screen.getByRole("textbox");
    await user.type(input, "tag1,");
    await user.type(input, "tag2,");
    await user.type(input, "tag3,");

    expect(screen.getByText("tag1")).toBeInTheDocument();
    expect(screen.getByText("tag2")).toBeInTheDocument();
    expect(screen.getByText("tag3")).toBeInTheDocument();
  });
});
