import { ActionFunctionArgs, redirect } from "@remix-run/node";
import { Form, useActionData } from "@remix-run/react";
import { AuthService } from "~/services/auth.service";

export const action = async ({ request }: ActionFunctionArgs) => {
  const formData = await request.formData();

  const { err, cookie } = await AuthService.login(formData);

  if (!err && cookie) {
    return redirect("/home", {
      headers: {
        "Set-Cookie": cookie
      }
    });
  }
  return { err };
}

export default function Login() {
  const data = useActionData<typeof action>()
  return <div>
    <Form method="POST">
      <input name="email" type="email" />
      <input name="password" type="password" />
      {data?.err}
      <button type="submit">Login</button>
    </Form>
  </div>
}
