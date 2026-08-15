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

const copyText = async (text) => {
  if (navigator.clipboard?.writeText && window.isSecureContext) {
    await navigator.clipboard.writeText(text);
    return;
  }
  const field = document.createElement("textarea");
  field.value = text;
  field.setAttribute("readonly", "");
  field.style.cssText = "position:fixed;opacity:0;pointer-events:none";
  document.body.append(field);
  field.select();
  const copied = document.execCommand("copy");
  field.remove();
  if (!copied) throw new Error("clipboard unavailable");
};

document.querySelectorAll(".copy-code").forEach((button) => {
  button.addEventListener("click", async () => {
    const code = button.closest(".code-panel")?.querySelector("code");
    if (!code) return;
    const label = button.getAttribute("aria-label") || "Copy to clipboard";
    button.disabled = true;
    try {
      await copyText(code.textContent);
      button.dataset.copyState = "copied";
      button.setAttribute("aria-label", "Copied");
      button.title = "Copied";
    } catch (_) {
      button.dataset.copyState = "error";
      button.setAttribute("aria-label", "Unable to copy");
      button.title = "Unable to copy";
    }
    window.setTimeout(() => {
      delete button.dataset.copyState;
      button.setAttribute("aria-label", label);
      button.title = "Copy to clipboard";
      button.disabled = false;
    }, 1600);
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
      if (!dragged || dragged === link || dragged.closest(".plugin-group") !== link.closest(".plugin-group")) return;
      event.preventDefault();
      link.classList.add("drag-over");
    });
    link.addEventListener("dragleave", () => link.classList.remove("drag-over"));
    link.addEventListener("drop", (event) => {
      if (!dragged || dragged === link || dragged.closest(".plugin-group") !== link.closest(".plugin-group")) return;
      event.preventDefault();
      const placeAfter = event.clientY > link.getBoundingClientRect().top + link.offsetHeight / 2;
      pluginNav.insertBefore(dragged, placeAfter ? link.nextSibling : link);
      link.classList.remove("drag-over");
      saveOrder();
    });
  });
}

const groupToggles = [...document.querySelectorAll(".plugin-group-toggle")];
if (groupToggles.length) {
  const saveGroups = async () => {
    const state = Object.fromEntries([...document.querySelectorAll(".plugin-group")].map((group) => [group.dataset.groupId, group.dataset.collapsed === "true"]));
    await fetch("/plugins/groups", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify(state) });
  };
  groupToggles.forEach((toggle) => toggle.addEventListener("click", () => {
    const group = toggle.closest(".plugin-group");
    const collapsed = group.dataset.collapsed !== "true";
    group.dataset.collapsed = String(collapsed);
    toggle.setAttribute("aria-expanded", String(!collapsed));
    saveGroups().catch(() => { group.dataset.collapsed = String(!collapsed); toggle.setAttribute("aria-expanded", String(collapsed)); });
  }));
}

const pluginSearchInput = document.querySelector("#plugin-search-input");
const pluginSearchData = document.querySelector("#plugin-search-data");
const pluginSearchResults = document.querySelector("#plugin-search-results");
if (pluginSearchInput && pluginSearchData && pluginSearchResults) {
  const plugins = JSON.parse(pluginSearchData.textContent);
  const matches = (query) => { const terms = query.toLocaleLowerCase().trim().split(/\s+/).filter(Boolean); if (!terms.length) return []; return plugins.filter((plugin) => { const haystack = [plugin.name, plugin.description, plugin.category, plugin.id, ...plugin.terms].join(" ").toLocaleLowerCase(); return terms.every((term) => haystack.includes(term)); }).slice(0, 7); };
  const render = () => { const results = matches(pluginSearchInput.value); pluginSearchResults.replaceChildren(); results.forEach((plugin) => { const link = document.createElement("a"); link.href = plugin.href; link.setAttribute("role", "option"); const name = document.createElement("strong"); name.textContent = plugin.name; const detail = document.createElement("span"); detail.textContent = plugin.description || plugin.category; link.append(name, detail); pluginSearchResults.append(link); }); if (pluginSearchInput.value.trim() && !results.length) { const empty = document.createElement("p"); empty.textContent = "No matching tools. Try a broader word."; pluginSearchResults.append(empty); } pluginSearchResults.hidden = !pluginSearchInput.value.trim(); };
  pluginSearchInput.addEventListener("input", render);
  pluginSearchInput.addEventListener("keydown", (event) => { if (event.key === "Escape") { pluginSearchInput.value = ""; pluginSearchResults.hidden = true; pluginSearchInput.blur(); } if (event.key === "Enter") { const first = pluginSearchResults.querySelector("a"); if (first) { event.preventDefault(); first.click(); } } });
  document.addEventListener("click", (event) => { if (!event.target.closest(".plugin-search")) pluginSearchResults.hidden = true; });
  document.addEventListener("keydown", (event) => { if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === "k") { event.preventDefault(); pluginSearchInput.focus(); } });
}
