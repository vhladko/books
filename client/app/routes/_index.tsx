import type { MetaFunction } from "@remix-run/node";

export const meta: MetaFunction = () => {
  return [
    { title: "Books" },
    { name: "description", content: "Books application" },
  ];
};

export default function Index() {
  return (
    <div>
      Hello books
    </div>
  );
}
