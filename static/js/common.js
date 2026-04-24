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



function renderComponentsHint(components) {
  const container = document.getElementById("components-hint");
  container.innerHTML = "";

  for (let i = 0; i < components.length; i++) {
    const row = document.createElement("div");
    row.classList.add("component-row");

    const img = document.createElement("img");
    img.src = components[i].imageUrl;
    img.alt = components[i].name;
    img.classList.add("component-icon");

    const span = document.createElement("span");
    span.textContent = components[i].name;

    row.appendChild(img);
    row.appendChild(span);
    container.appendChild(row);
  }
}


function createButtonRow() {
  let canGuess = true
  const secondChallenge = document.getElementById("2nd-challenge");
  const container = document.getElementById("Buttons");
   
  const introBox = document.createElement("div");
  introBox.classList.add("intro-box");

  const pageTitle = document.createElement("p");
  pageTitle.classList.add("page-title");
  pageTitle.textContent = "Bravo, maintenant lequel est-ce ?";

  introBox.append(pageTitle)

  container.innerHTML = "";
  ["Passive", "Q", "W", "E", "R"].forEach(letter => {
    const btn = document.createElement("button");
    if (letter != "Passive") {
      btn.textContent = letter;
    } else {
      btn.textContent = "P";
    }
     

    btn.addEventListener("click", () => {
          if (canGuess) { 
              compareSlot(letter, btn);
              canGuess = false 
              }
        });

    container.appendChild(btn);
  });

  secondChallenge.insertBefore(introBox, container);
}



function startDailyCountdown(challenge) {
          const timer = document.getElementById("dailyTimer");

          function updateTimer() {
              const now = new Date();

              const nextMidnight = new Date();
              nextMidnight.setHours(24, 0, 0, 0);

              const diff = nextMidnight - now;

              const hours = Math.floor(diff / 1000 / 60 / 60);
              const minutes = Math.floor((diff / 1000 / 60) % 60);
              const seconds = Math.floor((diff / 1000) % 60);

              timer.textContent = `Prochain ${challenge} dans ${hours}h ${minutes}m ${seconds}s`;
          }

          updateTimer();
          setInterval(updateTimer, 1000);
      }



function getTodayKey() {
          const d = new Date();
          return d.toISOString().slice(0,10);  
      }