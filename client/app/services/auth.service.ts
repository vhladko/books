class Service {
  async login(formData: FormData): Promise<{err: null | string, cookie: string | null}> {
    const response = await fetch('http://localhost:8080/api/login', { method: "POST", body: formData })
    const cookie = response.headers.get("Set-Cookie")

    console.log(cookie, response.headers.get("Set-Cookie"))
    if(response.ok) {
      return {err: null, cookie}
    }

    const json = await response.json();

    return {err: json.err, cookie: null};
  }
}

export const AuthService = new Service()

