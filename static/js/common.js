// Returns the elements whose name matches the current input value
function FilteredList(array, inputValue) {
  return array.filter(el =>
    el.name.toLowerCase().startsWith(inputValue.toLowerCase())
  );
}

// Displays clickable suggestions that have not been guessed yet
function renderList(array, list, guessed, input, submitGuess) {
  list.innerHTML = "";

  const displayedNames = new Set();

  for (let i = 0; i < array.length; i++) {
    if (guessed.includes(array[i].name)) continue;
    if (displayedNames.has(array[i].name)) continue;

    displayedNames.add(array[i].name);

    const item = document.createElement("div");

    const img = document.createElement("img");
    img.src = array[i].imageUrl;
    img.alt = array[i].name;

    const span = document.createElement("span");
    span.textContent = array[i].name;

    item.appendChild(img);
    item.appendChild(span);
    list.appendChild(item);

    item.addEventListener("click", function() {
      list.innerHTML = "";
      guessed.push(array[i].name);
      submitGuess(array[i].name);
      input.value = "";
    });
  }
}


const statConfig = {
  FlatHPPoolMod: {
    label: "HP",
    icon: "/static/images/stats/HP.png"
  },
  FlatMagicDamageMod: {
    label: "AP",
    icon: "/static/images/stats/AP.png"
  },
  FlatPhysicalDamageMod: {
    label: "AD",
    icon: "/static/images/stats/AD.png"
  },
  FlatArmorMod: {
    label: "Armor",
    icon: "/static/images/stats/Armor.png"
  },
  FlatSpellBlockMod: {
    label: "Magic Resist",
    icon: "/static/images/stats/MR.png"
  },
  PercentAttackSpeedMod: {
    label: "Attack Speed",
    icon: "/static/images/stats/AS.png"
  },
  FlatMovementSpeedMod: {
    label: "Move Speed",
    icon: "/static/images/stats/MS.png" 
  },
  FlatCritChanceMod: {
    label: "Crit Chance",
    icon: "/static/images/stats/Crit.png"
  },
  FlatMPPoolMod: {
    label: "Mana",
    icon: "/static/images/stats/Mana.png"
  }
};



function renderStats(stats) {
  const statsContainer = document.getElementById("stats");
  statsContainer.innerHTML = "";

  for (const key in stats) {
    const value = stats[key];

    if (!statConfig[key]) continue;

    const row = document.createElement("div");
    row.classList.add("stat-row");

    const img = document.createElement("img");
    img.src = statConfig[key].icon;
    img.alt = statConfig[key].label;
    img.classList.add("stat-icon");

    const span = document.createElement("span");
    span.textContent = `${statConfig[key].label} : ${value}`;

    row.appendChild(img);
    row.appendChild(span);
    statsContainer.appendChild(row);
  }
}