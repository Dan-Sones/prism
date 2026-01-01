import React from "react";
import type { ConfigurationElementType } from "./ConfigurationElement";
import TagBubble from "./TagBubble";

const Tag = (props: ConfigurationElementType) => {
  const [userInput, setUserInput] = React.useState<string>("");
  const [tags, setTags] = React.useState<Array<string>>([]);
  const [error, setError] = React.useState<string>("");

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setUserInput(e.target.value);
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    setError("");

    if (e.key === ",") {
      e.preventDefault();

      if (userInput.includes(" ")) {
        setError(
          "Tags cannot contain spaces. Please use hyphens or underscores."
        );
        return;
      }

      if (userInput && !tags.includes(userInput)) {
        setTags([...tags, userInput]);
      }
      setUserInput("");
    }
  };

  const removeTag = (tagToRemove: string) => {
    setTags(tags.filter((tag) => tag !== tagToRemove));
  };

  return (
    <div>
      <input
        type="text"
        name={props.name}
        value={userInput}
        onChange={handleInputChange}
        onKeyDown={handleKeyDown}
        className={`rounded bg-slate-50 h-9 p-3 w- text-slate-950 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-blue-500 ${
          error ? "border-2 border-red-500 focus:ring-red-500" : ""
        }`}
      />
      <p className="text-red-500 text-sm opacity-85 pt-1">{error}</p>
      <div className="mt-2 flex flex-wrap gap-2">
        {tags.map((tag, index) => (
          <TagBubble
            key={index}
            label={tag}
            onCrossClick={() => removeTag(tag)}
          />
        ))}
      </div>
    </div>
  );
};

export default Tag;
