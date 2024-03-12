import { ActionFunctionArgs, LoaderFunctionArgs } from "@remix-run/node";
import { json, useFetcher, useLoaderData } from "@remix-run/react"
import { BookService } from "~/services/book.service";
import { UserService } from "~/services/user.service";

export const loader = async ({ request }: LoaderFunctionArgs) => {
  const { books, err } = await UserService.getUserBooks({ request })

  return json({ books, err });
}

export const action = async ({ request }: ActionFunctionArgs) => {
  const formData = await request.formData();
  const { _action, ...values } = Object.fromEntries(formData);

  console.log(_action, "action inside");

  if (_action === "getBook") {
    const isbn = values.isbn as string;
    return BookService.getBookByIsbn({ isbn })
  }
}

export default function Home() {
  const { books, err } = useLoaderData<typeof loader>()
  const fetcher = useFetcher();
  console.log(fetcher.data, books, err);
  return <div>
    <h1>This is home</h1>
    <div>Your books</div>
    <div></div>

    <fetcher.Form method="post">
      <input type="search" name="isbn" />
      <button type="submit" name="_action" value="getBook">Find</button>
    </fetcher.Form>
  </div>
}
