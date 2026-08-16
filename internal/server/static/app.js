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

const sidebarVisibilityCheckboxes = [...document.querySelectorAll("[data-sidebar-visibility]")];
sidebarVisibilityCheckboxes.forEach((checkbox) => {
  if (checkbox.dataset.indeterminate === "true") checkbox.indeterminate = true;
  checkbox.addEventListener("change", () => checkbox.closest("form")?.requestSubmit());
});

const pluginSearchData = document.querySelector("#plugin-search-data");
const plugins = pluginSearchData ? JSON.parse(pluginSearchData.textContent) : [];
const pluginNavigationData = document.querySelector("#plugin-navigation-data");
const pluginNavigation = pluginNavigationData ? JSON.parse(pluginNavigationData.textContent) : { favorites: [], recents: [] };
const pluginByID = new Map(plugins.map((plugin) => [plugin.id, plugin]));

const initialPlugins = () => {
  const orderedIDs = [...pluginNavigation.recents, ...pluginNavigation.favorites];
  const seen = new Set();
  const pinned = orderedIDs.map((id) => pluginByID.get(id)).filter((plugin) => plugin && !seen.has(plugin.id) && seen.add(plugin.id));
  const remaining = plugins.filter((plugin) => !seen.has(plugin.id)).sort((left, right) => left.name.localeCompare(right.name));
  return [...pinned, ...remaining];
};

const searchPlugins = (query, limit = 8) => {
  const terms = query.toLocaleLowerCase().trim().split(/\s+/).filter(Boolean);
  if (!terms.length) return initialPlugins().slice(0, limit);
  return plugins.map((plugin) => {
    const name = plugin.name.toLocaleLowerCase();
    const description = plugin.description.toLocaleLowerCase();
    const category = plugin.category.toLocaleLowerCase();
    const id = plugin.id.toLocaleLowerCase();
    const aliases = plugin.terms.map((term) => term.toLocaleLowerCase());
    let score = 0;
    for (const term of terms) {
      if (name === term) score += 1000;
      else if (name.startsWith(term)) score += 800;
      else if (aliases.includes(term)) score += 700;
      else if (aliases.some((alias) => alias.startsWith(term))) score += 650;
      else if (name.includes(term)) score += 600;
      else if (id.includes(term)) score += 500;
      else if (aliases.some((alias) => alias.includes(term))) score += 450;
      else if (description.includes(term)) score += 300;
      else if (category.includes(term)) score += 200;
      else return null;
    }
    return { plugin, score };
  }).filter(Boolean).sort((left, right) => right.score - left.score || left.plugin.name.localeCompare(right.plugin.name)).slice(0, limit).map(({ plugin }) => plugin);
};

const pluginSearchInput = document.querySelector("#plugin-search-input");
const pluginSearchResults = document.querySelector("#plugin-search-results");
if (pluginSearchInput && pluginSearchData && pluginSearchResults) {
  const matches = (query) => query.trim() ? searchPlugins(query, 7) : [];
  const render = () => { const results = matches(pluginSearchInput.value); pluginSearchResults.replaceChildren(); results.forEach((plugin) => { const link = document.createElement("a"); link.href = plugin.href; link.setAttribute("role", "option"); const name = document.createElement("strong"); name.textContent = plugin.name; const detail = document.createElement("span"); detail.textContent = plugin.description || plugin.category; link.append(name, detail); pluginSearchResults.append(link); }); if (pluginSearchInput.value.trim() && !results.length) { const empty = document.createElement("p"); empty.textContent = "No matching tools. Try a broader word."; pluginSearchResults.append(empty); } pluginSearchResults.hidden = !pluginSearchInput.value.trim(); };
  pluginSearchInput.addEventListener("input", render);
  pluginSearchInput.addEventListener("keydown", (event) => { if (event.key === "Escape") { pluginSearchInput.value = ""; pluginSearchResults.hidden = true; pluginSearchInput.blur(); } if (event.key === "Enter") { const first = pluginSearchResults.querySelector("a"); if (first) { event.preventDefault(); first.click(); } } });
  document.addEventListener("click", (event) => { if (!event.target.closest(".plugin-search")) pluginSearchResults.hidden = true; });
}

const commandPalette = document.querySelector("#command-palette");
const commandPaletteInput = document.querySelector("#command-palette-input");
const commandPaletteResults = document.querySelector("#command-palette-results");
const commandPaletteTrigger = document.querySelector("#command-palette-trigger");
if (commandPalette && commandPaletteInput && commandPaletteResults) {
  let activeIndex = 0;
  let lastFocusedElement = null;
  let currentResults = [];

  const renderCommandPalette = () => {
    currentResults = searchPlugins(commandPaletteInput.value);
    activeIndex = Math.min(activeIndex, Math.max(0, currentResults.length - 1));
    commandPaletteResults.replaceChildren();
    if (!currentResults.length) {
      const empty = document.createElement("p");
      empty.className = "command-palette-empty";
      empty.textContent = "No installed tools match that search.";
      commandPaletteResults.append(empty);
      commandPaletteInput.removeAttribute("aria-activedescendant");
      return;
    }
    currentResults.forEach((plugin, index) => {
      const link = document.createElement("a");
      link.id = `command-palette-result-${index}`;
      link.href = plugin.href;
      link.setAttribute("role", "option");
      link.setAttribute("aria-selected", String(index === activeIndex));
      link.className = index === activeIndex ? "is-active" : "";
      const name = document.createElement("strong");
      name.textContent = plugin.name;
      const detail = document.createElement("span");
      detail.textContent = plugin.description || plugin.category;
      const category = document.createElement("em");
      category.textContent = plugin.category || "tool";
      link.append(name, detail, category);
      commandPaletteResults.append(link);
    });
    commandPaletteInput.setAttribute("aria-activedescendant", `command-palette-result-${activeIndex}`);
  };

  const closeCommandPalette = () => {
    if (commandPalette.hidden) return;
    commandPalette.hidden = true;
    document.body.classList.remove("command-palette-open");
    commandPaletteInput.value = "";
    if (lastFocusedElement instanceof HTMLElement) lastFocusedElement.focus();
  };

  const openCommandPalette = () => {
    if (!commandPalette.hidden) return;
    lastFocusedElement = document.activeElement;
    commandPalette.hidden = false;
    document.body.classList.add("command-palette-open");
    activeIndex = 0;
    renderCommandPalette();
    commandPaletteInput.focus();
  };

  commandPaletteTrigger?.addEventListener("click", openCommandPalette);
  commandPalette.querySelectorAll("[data-command-palette-close]").forEach((button) => button.addEventListener("click", closeCommandPalette));
  commandPaletteInput.addEventListener("input", () => { activeIndex = 0; renderCommandPalette(); });
  commandPaletteInput.addEventListener("keydown", (event) => {
    if (event.key === "Escape") { event.preventDefault(); closeCommandPalette(); return; }
    if (event.key === "ArrowDown" || event.key === "ArrowUp") {
      event.preventDefault();
      if (!currentResults.length) return;
      activeIndex = (activeIndex + (event.key === "ArrowDown" ? 1 : -1) + currentResults.length) % currentResults.length;
      renderCommandPalette();
      return;
    }
    if (event.key === "Enter" && currentResults[activeIndex]) {
      event.preventDefault();
      window.location.assign(currentResults[activeIndex].href);
    }
  });
  document.addEventListener("keydown", (event) => {
    if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === "k") {
      event.preventDefault();
      openCommandPalette();
    } else if (event.key === "Escape" && !commandPalette.hidden) {
      closeCommandPalette();
    }
  });
}
