async function fetchData() {
  try {
    const res = await fetch("http://localhost:3000/");
    const data = res.json();
    console.log(res);
  } catch (error) {
    console.log(error);
  }
}

fetchData();
