interface Book { }

class Service {
  async getBookByIsbn({ isbn }: { isbn: string }): Promise<{ err: null | string, book: Book | null }> {
    const response = await fetch('http://localhost:8080/api/book/isbn/' + isbn)

    console.log('afeter fetch');

    const data = await response.json()
    console.log(data)

    if (response.ok) {
      return { err: null, book: data.book }
    }


    return { err: data.err, book: [] };
  }
}

export const BookService = new Service()
