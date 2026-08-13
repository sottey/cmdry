document.querySelectorAll("table[data-sortable]").forEach((table) => {
  table.querySelectorAll("th").forEach((header, index) => {
    header.addEventListener("click", () => {
      const body = table.tBodies[0];
      const ascending = header.dataset.direction !== "asc";
      table.querySelectorAll("th").forEach((cell) => { cell.dataset.direction = ""; });
      header.dataset.direction = ascending ? "asc" : "desc";
      [...body.rows].sort((left, right) => {
        const a = left.cells[index].textContent.trim();
        const b = right.cells[index].textContent.trim();
        const numeric = Number(a) - Number(b);
        const comparison = Number.isNaN(numeric) ? a.localeCompare(b) : numeric;
        return ascending ? comparison : -comparison;
      }).forEach((row) => body.appendChild(row));
    });
  });
});
