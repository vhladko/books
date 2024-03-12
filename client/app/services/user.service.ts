class Service {
  async getUserBooks({ request }: { request: Request }): Promise<{ err: null | string, books: unknown[] }> {
    const response = await fetch('http://localhost:8080/api/books')
    const data = await response.json()

    if (response.ok) {
      return { err: null, books: data.books }
    }


    return { err: data.err, books: [] };
  }
}

export const UserService = new Service()
