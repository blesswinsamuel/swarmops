export default async function fetcher(url) {
  const res = await fetch(url);
  if (res.ok) {
    return await res.json();
  }
  throw new Error("Failed to fetch");
}
