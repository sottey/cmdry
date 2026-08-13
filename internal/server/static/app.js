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

const pluginNav = document.querySelector("#plugin-nav");
if (pluginNav) {
  let dragged = null;
  const pluginLinks = () => [...pluginNav.querySelectorAll("a[data-plugin-id]")];
  const saveOrder = async () => {
    const ids = pluginLinks().map((link) => link.dataset.pluginId);
    try {
      const response = await fetch("/plugins/order", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(ids),
      });
      if (!response.ok) throw new Error("save failed");
    } catch (_) {
      pluginNav.classList.add("plugin-order-error");
    }
  };
  pluginLinks().forEach((link) => {
    link.addEventListener("dragstart", (event) => {
      dragged = link;
      event.dataTransfer.effectAllowed = "move";
      event.dataTransfer.setData("text/plain", link.dataset.pluginId);
      link.classList.add("dragging");
    });
    link.addEventListener("dragend", () => {
      dragged?.classList.remove("dragging");
      dragged = null;
      pluginNav.querySelectorAll(".drag-over").forEach((item) => item.classList.remove("drag-over"));
    });
    link.addEventListener("dragover", (event) => {
      if (!dragged || dragged === link) return;
      event.preventDefault();
      link.classList.add("drag-over");
    });
    link.addEventListener("dragleave", () => link.classList.remove("drag-over"));
    link.addEventListener("drop", (event) => {
      if (!dragged || dragged === link) return;
      event.preventDefault();
      const placeAfter = event.clientY > link.getBoundingClientRect().top + link.offsetHeight / 2;
      pluginNav.insertBefore(dragged, placeAfter ? link.nextSibling : link);
      link.classList.remove("drag-over");
      saveOrder();
    });
  });
}
