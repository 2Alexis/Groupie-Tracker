package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var (
	favoriteWeapons  []Weapon
	favoriteMonsters []Monster
	favoriteObjets   []Objets
	favoriteSkills   []Skills
	favoriteEvents   []Evenements
	favoriteDeco     []Decorations
	favoriteCharms   []Charms
	favoriteArmors   []Armor
)

type Favorites struct {
	Weapons  []Weapon  // Liste des armes favorites
	Monsters []Monster // Liste des monstres favoris
	Objets   []Objets  // Liste des armes favorites
	Skills   []Skills
	Events   []Evenements
	Deco     []Decorations
	Charms   []Charms
	Armors   []Armor
}

type FavoritesItems struct {
	ID       int    `json:"id"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Rarete   int    `json:"rarete"`
	Category string `json:"category"`
	// Autres champs nécessaires pour représenter un élément
}

var favoriteItems []FavoritesItems

// Structure pour stocker les données des armes

type PageData struct {
	Weapons  []Weapon
	Monsters []Monster
	Objets   []Objets
	Skills   []Skills
	Events   []Evenements
	Deco     []Decorations
	Charms   []Charms
	Armors   []Armor
}

// Structure pour stocker les données des armes
type Weapon struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Type     string `json:"type"`
	Rarity   int    `json:"rarity"`
	Attack   struct {
		Display int `json:"display"`
		Raw     int `json:"raw"`
	} `json:"attack"`
	Element []struct {
		Type   string `json:"type"`
		Damage int    `json:"damage"`
	} `json:"elements"`

	Assets struct {
		Image string `json:"image"`
	} `json:"assets"`
	Rang        string `json:"rank"`
	TypeDegats  string `json:"damageType"`
	Resistances []struct {
		Element string `json:"element"`
	} `json:"resistances"`
	QRang       int    `json:"questRank"`
	Description string `json:"description"`
	Size        string `json:"type"`
}

type Skills struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Rarity      int    `json:"rarity"`
	Rank        []struct {
		Niveau      int    `json:"level"`
		Description string `json:"description"`
	} `json:"ranks"`
	Skills []struct {
		Name        string `json:"skillName"`
		Niveau      int    `json:"level"`
		Description string `json:"description"`
	} `json:"skills"`
	Element []struct {
		Type   string `json:"type"`
		Damage int    `json:"damage"`
	} `json:"elements"`
	TypeDegats  string `json:"damageType"`
	Rang        string `json:"rank"`
	Resistances []struct {
		Element string `json:"element"`
	} `json:"resistances"`
	QRang int    `json:"questRank"`
	Size  string `json:"type"`
}

type Objets struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Rarity      int    `json:"rarity"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Type        string `json:"type"`
	Rank        []struct {
		Niveau      int    `json:"level"`
		Description string `json:"description"`
	} `json:"ranks"`
	Skills []struct {
		Name        string `json:"skillName"`
		Niveau      int    `json:"level"`
		Description string `json:"description"`
	} `json:"skills"`
	Rang        string `json:"rank"`
	Resistances []struct {
		Element string `json:"element"`
	} `json:"resistances"`
	QRang int `json:"questRank"`

	Size    string `json:"type"`
	Element []struct {
		Type   string `json:"type"`
		Damage int    `json:"damage"`
	} `json:"elements"`
	TypeDegats string `json:"damageType"`
}

type Decorations struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Rarity   int    `json:"rarity"`
	Type     string `json:"type"`
	Skills   []struct {
		Name        string `json:"skillName"`
		Niveau      int    `json:"level"`
		Description string `json:"description"`
	} `json:"skills"`
	Rank []struct {
		Niveau      int    `json:"level"`
		Description string `json:"description"`
	} `json:"ranks"`
	Rang        string `json:"rank"`
	Resistances []struct {
		Element string `json:"element"`
	} `json:"resistances"`
	QRang       int    `json:"questRank"`
	Description string `json:"description"`
	Size        string `json:"type"`
	Element     []struct {
		Type   string `json:"type"`
		Damage int    `json:"damage"`
	} `json:"elements"`
	TypeDegats string `json:"damageType"`
}

type Charms struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Type     string `json:"type"`
	Rarity   int    `json:"rarity"`
	CRank    []struct {
		Niveau int `json:"level"`
		Rarete int `json:"rarity"`
		Skills []struct {
			Name        string `json:"skillName"`
			Niveau      int    `json:"level"`
			Description string `json:"description"`
		}
	} `json:"ranks"`
	Rang        string `json:"rank"`
	Resistances []struct {
		Element string `json:"element"`
	} `json:"resistances"`
	QRang       int    `json:"questRank"`
	Description string `json:"description"`
	Size        string `json:"type"`
	Element     []struct {
		Type   string `json:"type"`
		Damage int    `json:"damage"`
	} `json:"elements"`
	TypeDegats string `json:"damageType"`
}

type Armor struct {
	ID          int            `json:"id"`
	Name        string         `json:"name"`
	Category    string         `json:"category"`
	Type        string         `json:"type"`
	Rarity      int            `json:"rarity"`
	Rang        string         `json:"rank"`
	Resistances map[string]int `json:"resistances"`
	Skills      []struct {
		Name        string `json:"skillName"`
		Niveau      int    `json:"level"`
		Description string `json:"description"`
	} `json:"skills"`

	Defense struct {
		Base    int `json:"base"`
		Max     int `json:"max"`
		Upgrade int `json:"augmented"`
	} `json:"defense"`
	Assets struct {
		ImageH string `json:"imageMale"`
		ImageF string `json:"imageFemale"`
	} `json:"assets"`

	QRang       int    `json:"questRank"`
	Description string `json:"description"`
	Size        string `json:"type"`
	Element     []struct {
		Type   string `json:"type"`
		Damage int    `json:"damage"`
	} `json:"elements"`
	TypeDegats string `json:"damageType"`
}

// Structure pour stocker les données des monstres
type Monster struct {
	ID          int    `json:"id"`
	Category    string `json:"category"` // Type de monstre
	Name        string `json:"name"`
	Species     string `json:"species"`
	Size        string `json:"type"`
	Type        string `json:"type"`
	Rarity      int    `json:"rarity"`
	Description string `json:"description"`
	Recompense  []struct {
		Item struct {
			Name string `json:"name"`
		} `json:"item"`
	} `json:"rewards"`
	Weaknesses []struct {
		Element string `json:"element"`
		Stars   int    `json:"stars"`
	} `json:"weaknesses"`
	Ailments []struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"ailments"`
	Resistances []struct {
		Element string `json:"element"`
	} `json:"resistances"`
	Rang  string `json:"rank"`
	QRang int    `json:"questRank"`

	TypeDegats string `json:"damageType"`
}

type Evenements struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Category     string `json:"category"`
	Type         string `json:"type"`
	Rarity       int    `json:"rarity"`
	QRang        int    `json:"questRank"`
	Description  string `json:"description"`
	Exigences    string `json:"requirements"`
	Objectif     string `json:"successConditions"`
	Localisation struct {
		Name string `json:"name"`
	} `json:"location"`
	Rang        string `json:"rank"`
	Resistances []struct {
		Element string `json:"element"`
	} `json:"resistances"`
	Size    string `json:"type"`
	Element []struct {
		Type   string `json:"type"`
		Damage int    `json:"damage"`
	} `json:"elements"`
	TypeDegats string `json:"damageType"`
}

var armors []Armor

// Fonction pour récupérer les données à partir de l'API
func fetchJSONData(url string, data interface{}) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(data)
	if err != nil {
		return err
	}
	return nil
}
func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	// Obtenir les paramètres de filtrage de la requête
	typeFilterEvent := r.URL.Query().Get("typeFilterEvent")
	typeFilterArmor := r.URL.Query().Get("typeFilterArmor")
	rankFilter := r.URL.Query().Get("rankFilter")

	rankFilterArmor := r.URL.Query().Get("rankFilterArmor")
	rarityFilter := r.URL.Query().Get("rarityFilter")
	elementFilter := r.URL.Query().Get("elementFilter")
	damageTypeFilter := r.URL.Query().Get("damageTypeFilter")
	sizeFilter := r.URL.Query().Get("sizeFilter")
	resistanceFilter := r.URL.Query().Get("resistanceFilter")
	query := r.URL.Query().Get("q")

	var weapons []Weapon
	weaponsURL := "https://mhw-db.com/weapons"
	err := fetchJSONData(weaponsURL, &weapons)
	if err != nil {
		log.Fatal(err)
	}

	var filteredWeapons []Weapon
	for _, weapon := range weapons {
		// Vérifiez si l'arme correspond aux critères de filtrage
		if (sizeFilter == "" || weapon.Size == sizeFilter) &&
			(resistanceFilter == "" || containsResistance(weapon.Resistances, resistanceFilter)) &&
			(typeFilterEvent == "" || weapon.Type == typeFilterEvent) &&
			(rankFilter == "" || strconv.Itoa(weapon.QRang) == rankFilter) &&
			(typeFilterArmor == "" || weapon.Type == typeFilterArmor) &&
			(rankFilterArmor == "" || weapon.Rang == rankFilterArmor) &&
			(rarityFilter == "" || strconv.Itoa(weapon.Rarity) == rarityFilter) &&
			(elementFilter == "" || containsElement(weapon.Element, elementFilter)) &&
			(damageTypeFilter == "" || weapon.TypeDegats == damageTypeFilter) &&
			(query == "" || strings.Contains(strings.ToLower(weapon.Name), strings.ToLower(query))) {
			filteredWeapons = append(filteredWeapons, weapon)
		}
	}

	// Filtrer les armures en fonction des critères de filtrage

	var armors []Armor
	armorsURL := "https://mhw-db.com/armor"
	errArmor := fetchJSONData(armorsURL, &armors)
	if errArmor != nil {
		log.Fatal(errArmor)
	}
	if len(armors) == 0 {
		log.Println("Aucune donnée n'a été reçue pour les armes.")
	} else {
		log.Println("Données reçues pour les armes:", weapons)
	}

	var filteredArmors []Armor
	for _, armor := range armors {
		// Vérifiez si l'armure correspond aux critères de filtrage
		if (sizeFilter == "" || armor.Size == sizeFilter) &&
			(resistanceFilter == "" || armor.Resistances[resistanceFilter] > 0) &&
			(typeFilterEvent == "" || armor.Type == typeFilterEvent) &&
			(rankFilter == "" || strconv.Itoa(armor.QRang) == rankFilter) &&
			(typeFilterArmor == "" || armor.Type == typeFilterArmor) &&
			(rankFilterArmor == "" || armor.Rang == rankFilterArmor) &&
			(rarityFilter == "" || strconv.Itoa(armor.Rarity) == rarityFilter) &&
			(elementFilter == "" || containsElement(armor.Element, elementFilter)) &&
			(damageTypeFilter == "" || armor.TypeDegats == damageTypeFilter) &&

			(query == "" || strings.Contains(strings.ToLower(armor.Name), strings.ToLower(query))) {
			filteredArmors = append(filteredArmors, armor)
		}
	}

	var charms []Charms
	charmsURL := "https://mhw-db.com/charms"
	errcharms := fetchJSONData(charmsURL, &charms)
	if errcharms != nil {
		log.Fatal(errcharms)
	}

	var filteredCharms []Charms
	for _, charm := range charms {
		if (sizeFilter == "" || charm.Size == sizeFilter) &&
			(resistanceFilter == "" || containsResistance(charm.Resistances, resistanceFilter)) &&
			(typeFilterEvent == "" || charm.Type == typeFilterEvent) &&
			(rankFilter == "" || strconv.Itoa(charm.QRang) == rankFilter) &&
			(typeFilterArmor == "" || charm.Type == typeFilterArmor) &&
			(rankFilterArmor == "" || charm.Rang == rankFilterArmor) &&
			(rarityFilter == "" || strconv.Itoa(charm.Rarity) == rarityFilter) &&
			(elementFilter == "" || containsElement(charm.Element, elementFilter)) &&
			(damageTypeFilter == "" || charm.TypeDegats == damageTypeFilter) &&
			query == "" || strings.Contains(strings.ToLower(charm.Name), strings.ToLower(query)) {
			filteredCharms = append(filteredCharms, charm)
		}
	}

	var deco []Decorations
	decoURL := "https://mhw-db.com/decorations"
	errDeco := fetchJSONData(decoURL, &deco)
	if errDeco != nil {
		log.Fatal(errDeco)
	}

	// Filtrer les décorations en fonction des critères de filtrage
	var filteredDeco []Decorations
	for _, decoration := range deco {
		// Vérifiez si la décoration correspond aux critères de filtrage
		if (sizeFilter == "" || decoration.Size == sizeFilter) &&
			(resistanceFilter == "" || containsResistance(decoration.Resistances, resistanceFilter)) &&
			(typeFilterEvent == "" || decoration.Type == typeFilterEvent) &&
			(rankFilter == "" || strconv.Itoa(decoration.QRang) == rankFilter) &&
			(typeFilterArmor == "" || decoration.Type == typeFilterArmor) &&
			(rankFilterArmor == "" || decoration.Rang == rankFilterArmor) &&
			(rarityFilter == "" || strconv.Itoa(decoration.Rarity) == rarityFilter) &&
			(elementFilter == "" || containsElement(decoration.Element, elementFilter)) &&
			(damageTypeFilter == "" || decoration.TypeDegats == damageTypeFilter) &&
			(query == "" || strings.Contains(strings.ToLower(decoration.Name), strings.ToLower(query))) {
			filteredDeco = append(filteredDeco, decoration)
		}
	}

	var events []Evenements
	eventsURL := "https://mhw-db.com/events"
	errevents := fetchJSONData(eventsURL, &events)
	if errevents != nil {
		log.Fatal(errevents)
	}

	var filteredEvents []Evenements
	for _, event := range events {
		// Vérifiez si l'événement correspond aux critères de filtrage
		if (sizeFilter == "" || event.Size == sizeFilter) &&
			(resistanceFilter == "" || containsResistance(event.Resistances, resistanceFilter)) &&
			(typeFilterEvent == "" || event.Type == typeFilterEvent) &&
			(rankFilter == "" || strconv.Itoa(event.QRang) == rankFilter) &&
			(typeFilterArmor == "" || event.Type == typeFilterArmor) &&
			(rankFilterArmor == "" || event.Rang == rankFilterArmor) &&
			(rarityFilter == "" || strconv.Itoa(event.Rarity) == rarityFilter) &&
			(elementFilter == "" || containsElement(event.Element, elementFilter)) &&
			(damageTypeFilter == "" || event.TypeDegats == damageTypeFilter) &&
			(query == "" || strings.Contains(strings.ToLower(event.Name), strings.ToLower(query))) {
			filteredEvents = append(filteredEvents, event)
		}
	}

	var items []Objets
	itemsURL := "https://mhw-db.com/items"
	errItems := fetchJSONData(itemsURL, &items)
	if errItems != nil {
		log.Fatal(errItems)
	}

	var filteredItems []Objets
	for _, item := range items {
		// Vérifiez si l'objet correspond aux critères de filtrage
		if (sizeFilter == "" || item.Size == sizeFilter) &&
			(resistanceFilter == "" || containsResistance(item.Resistances, resistanceFilter)) &&
			(typeFilterEvent == "" || item.Type == typeFilterEvent) &&
			(rankFilter == "" || strconv.Itoa(item.QRang) == rankFilter) &&
			(typeFilterArmor == "" || item.Type == typeFilterArmor) &&
			(rankFilterArmor == "" || item.Rang == rankFilterArmor) &&
			(rarityFilter == "" || strconv.Itoa(item.Rarity) == rarityFilter) &&
			(elementFilter == "" || containsElement(item.Element, elementFilter)) &&
			(damageTypeFilter == "" || item.TypeDegats == damageTypeFilter) &&
			(query == "" || strings.Contains(strings.ToLower(item.Name), strings.ToLower(query))) {
			filteredItems = append(filteredItems, item)
		}
	}

	var monsters []Monster
	monstersURL := "https://mhw-db.com/monsters"
	errmonsters := fetchJSONData(monstersURL, &monsters)
	if errmonsters != nil {
		log.Fatal(errmonsters)
	}

	var filteredMonsters []Monster
	for _, monster := range monsters {
		// Vérifiez si le monstre correspond aux critères de filtrage
		if (sizeFilter == "" || monster.Size == sizeFilter) &&
			(resistanceFilter == "" || containsResistance(monster.Resistances, resistanceFilter)) &&
			(typeFilterEvent == "" || monster.Type == typeFilterEvent) &&
			(rankFilter == "" || strconv.Itoa(monster.QRang) == rankFilter) &&
			(typeFilterArmor == "" || monster.Type == typeFilterArmor) &&
			(rankFilterArmor == "" || monster.Rang == rankFilterArmor) &&
			(rarityFilter == "" || strconv.Itoa(monster.Rarity) == rarityFilter) &&

			(damageTypeFilter == "" || monster.TypeDegats == damageTypeFilter) &&
			(query == "" || strings.Contains(strings.ToLower(monster.Name), strings.ToLower(query))) {
			filteredMonsters = append(filteredMonsters, monster)
		}
	}

	var skills []Skills
	skillsURL := "https://mhw-db.com/skills"
	errskills := fetchJSONData(skillsURL, &skills)
	if errskills != nil {
		log.Fatal(errskills)
	}

	var filteredSkills []Skills
	for _, skill := range skills {
		if (sizeFilter == "" || skill.Size == sizeFilter) &&
			(resistanceFilter == "" || containsResistance(skill.Resistances, resistanceFilter)) &&
			(typeFilterEvent == "" || skill.Type == typeFilterEvent) &&
			(rankFilter == "" || strconv.Itoa(skill.QRang) == rankFilter) &&
			(typeFilterArmor == "" || skill.Type == typeFilterArmor) &&
			(rankFilterArmor == "" || skill.Rang == rankFilterArmor) &&
			(rarityFilter == "" || strconv.Itoa(skill.Rarity) == rarityFilter) &&
			(elementFilter == "" || containsElement(skill.Element, elementFilter)) &&
			(damageTypeFilter == "" || skill.TypeDegats == damageTypeFilter) &&
			query == "" || strings.Contains(strings.ToLower(skill.Name), strings.ToLower(query)) {
			filteredSkills = append(filteredSkills, skill)
		}
	}

	data := struct {
		Weapons  []Weapon
		Armors   []Armor
		Charms   []Charms
		Deco     []Decorations
		Events   []Evenements
		Items    []Objets
		Skills   []Skills
		Monsters []Monster
		Query    string
	}{
		Weapons:  filteredWeapons,
		Armors:   filteredArmors,
		Charms:   filteredCharms,
		Deco:     filteredDeco,
		Events:   filteredEvents,
		Items:    filteredItems,
		Skills:   filteredSkills,
		Monsters: filteredMonsters,
		Query:    query,
	}
	// Autres traitements...

	// Charger le modèle HTML pour afficher les armes
	t, errTemplate := template.ParseFiles("templates/search.html")
	if errTemplate != nil {
		log.Fatal(errTemplate)
	}

	// Rendre la page HTML avec les données récupérées
	errExecute := t.Execute(w, data)
	if errExecute != nil {
		log.Fatal(errExecute)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Chargez le fichier HTML de la page d'accueil
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal(err)
	}
	// Exécutez le modèle pour générer la sortie HTML
	err = t.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func skillsHandler(w http.ResponseWriter, r *http.Request) {
	var skills []Skills
	skillsURL := "https://mhw-db.com/skills"
	err := fetchJSONData(skillsURL, &skills)
	if err != nil {
		log.Fatal(err)
	}

	// Obtenir le paramètre de recherche de la requête
	query := r.URL.Query().Get("q")

	// Filtrer les compétences en fonction de la recherche
	var filteredSkills []Skills
	for _, skill := range skills {
		if query == "" || strings.Contains(strings.ToLower(skill.Name), strings.ToLower(query)) {
			filteredSkills = append(filteredSkills, skill)
		}
	}

	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1 // Si le numéro de page est invalide ou non défini, affichez la première page par défaut
	}

	// Calculez l'indice de début et de fin pour paginer les compétences
	startIndex := (page - 1) * 10
	endIndex := startIndex + 10
	if endIndex > len(filteredSkills) {
		endIndex = len(filteredSkills)
	}

	// Sélectionnez les compétences à afficher sur cette page
	paginatedSkills := filteredSkills[startIndex:endIndex]

	// Préparez les données nécessaires pour la pagination
	hasNextPage := endIndex < len(filteredSkills)
	hasPreviousPage := page > 1
	nextPage := page + 1
	previousPage := page - 1

	// Chargez le modèle HTML avec les données paginées et les informations de pagination
	data := struct {
		Skills          []Skills
		HasNextPage     bool
		HasPreviousPage bool
		NextPage        int
		PreviousPage    int
		Query           string
	}{
		Skills:          paginatedSkills,
		HasNextPage:     hasNextPage,
		HasPreviousPage: hasPreviousPage,
		NextPage:        nextPage,
		PreviousPage:    previousPage,
		Query:           query,
	}

	t, err := template.ParseFiles("templates/skills.html")
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}
func skillDetailsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the weapon ID from the query parameters
	skillID := r.URL.Query().Get("id")
	if skillID == "" {
		http.Error(w, "Missing charm ID", http.StatusBadRequest)
		return
	}

	// Fetch details of the weapon with the given ID
	skillURL := fmt.Sprintf("https://mhw-db.com/skills/%s", skillID)
	var skill Skills
	err := fetchJSONData(skillURL, &skill)
	if err != nil {
		http.Error(w, "Failed to fetch charm details", http.StatusInternalServerError)
		log.Println("Failed to fetch charm details:", err)
		return
	}

	// Parse the HTML template for displaying weapon details
	t, err := template.ParseFiles("templates/skill_details.html")
	if err != nil {
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		log.Println("Failed to parse template:", err)
		return
	}

	// Render the template with the weapon details
	err = t.Execute(w, skill)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		log.Println("Failed to render template:", err)
		return
	}
}

func itemsHandler(w http.ResponseWriter, r *http.Request) {
	var items []Objets
	itemsURL := "https://mhw-db.com/items"
	err := fetchJSONData(itemsURL, &items)
	if err != nil {
		log.Fatal(err)
	}

	// Obtenir les paramètres de filtrage de la requête
	rarityFilter := r.URL.Query().Get("rarityFilter")
	searchQuery := r.URL.Query().Get("q")

	// Filtrer les objets en fonction des critères de filtrage
	var filteredItems []Objets
	for _, item := range items {
		// Vérifiez si l'objet correspond aux critères de filtrage
		if (rarityFilter == "" || strconv.Itoa(item.Rarity) == rarityFilter) &&
			(searchQuery == "" || strings.Contains(strings.ToLower(item.Name), strings.ToLower(searchQuery))) {
			filteredItems = append(filteredItems, item)
		}
	}

	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1 // Si le numéro de page est invalide ou non défini, affichez la première page par défaut
	}

	// Calculez l'indice de début et de fin pour paginer les objets
	startIndex := (page - 1) * 10
	endIndex := startIndex + 10
	if endIndex > len(filteredItems) {
		endIndex = len(filteredItems)
	}

	// Sélectionnez les objets à afficher sur cette page
	paginatedItems := filteredItems[startIndex:endIndex]

	// Préparez les données nécessaires pour la pagination
	hasNextPage := endIndex < len(filteredItems)
	hasPreviousPage := page > 1
	nextPage := page + 1
	previousPage := page - 1

	// Chargez le modèle HTML avec les données paginées et les informations de pagination
	data := struct {
		Items           []Objets
		HasNextPage     bool
		HasPreviousPage bool
		NextPage        int
		PreviousPage    int
		Query           string
		RarityFilter    string
	}{
		Items:           paginatedItems,
		HasNextPage:     hasNextPage,
		HasPreviousPage: hasPreviousPage,
		NextPage:        nextPage,
		PreviousPage:    previousPage,
		Query:           searchQuery,
		RarityFilter:    rarityFilter,
	}

	t, err := template.ParseFiles("templates/items.html")
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}
func itemDetailsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the weapon ID from the query parameters
	itemID := r.URL.Query().Get("id")
	if itemID == "" {
		http.Error(w, "Missing charm ID", http.StatusBadRequest)
		return
	}

	// Fetch details of the weapon with the given ID
	itemURL := fmt.Sprintf("https://mhw-db.com/items/%s", itemID)
	var item Objets
	err := fetchJSONData(itemURL, &item)
	if err != nil {
		http.Error(w, "Failed to fetch charm details", http.StatusInternalServerError)
		log.Println("Failed to fetch charm details:", err)
		return
	}

	// Parse the HTML template for displaying weapon details
	t, err := template.ParseFiles("templates/item_details.html")
	if err != nil {
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		log.Println("Failed to parse template:", err)
		return
	}

	// Render the template with the weapon details
	err = t.Execute(w, item)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		log.Println("Failed to render template:", err)
		return
	}
}

func eventsHandler(w http.ResponseWriter, r *http.Request) {
	var events []Evenements
	eventsURL := "https://mhw-db.com/events"
	err := fetchJSONData(eventsURL, &events)
	if err != nil {
		log.Fatal(err)
	}
	// Obtenir les paramètres de filtrage de la requête
	typeFilter := r.URL.Query().Get("typeFilter")
	rankFilter := r.URL.Query().Get("rankFilter")
	query := r.URL.Query().Get("q")

	// Filtrer les événements en fonction des critères de filtrage
	var filteredEvents []Evenements
	for _, event := range events {
		// Vérifiez si l'événement correspond aux critères de filtrage
		if (typeFilter == "" || event.Type == typeFilter) &&
			(rankFilter == "" || strconv.Itoa(event.QRang) == rankFilter) &&
			(query == "" || strings.Contains(strings.ToLower(event.Name), strings.ToLower(query))) {
			filteredEvents = append(filteredEvents, event)
		}
	}

	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1 // Si le numéro de page est invalide ou non défini, affichez la première page par défaut
	}

	// Calculez l'indice de début et de fin pour paginer les événements
	startIndex := (page - 1) * 10
	endIndex := startIndex + 10
	if endIndex > len(filteredEvents) {
		endIndex = len(filteredEvents)
	}

	// Sélectionnez les événements à afficher sur cette page
	paginatedEvents := filteredEvents[startIndex:endIndex]

	// Préparez les données nécessaires pour la pagination
	hasNextPage := endIndex < len(filteredEvents)
	hasPreviousPage := page > 1
	nextPage := page + 1
	previousPage := page - 1

	// Chargez le modèle HTML avec les données paginées et les informations de pagination
	data := struct {
		Events          []Evenements
		HasNextPage     bool
		HasPreviousPage bool
		NextPage        int
		PreviousPage    int
		Query           string
		TypeFilter      string
		RankFilter      string
	}{
		Events:          paginatedEvents,
		HasNextPage:     hasNextPage,
		HasPreviousPage: hasPreviousPage,
		NextPage:        nextPage,
		PreviousPage:    previousPage,
		Query:           query,
		TypeFilter:      typeFilter,
		RankFilter:      rankFilter,
	}

	t, err := template.ParseFiles("templates/events.html")
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}

func eventDetailsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the weapon ID from the query parameters
	eventID := r.URL.Query().Get("id")
	if eventID == "" {
		http.Error(w, "Missing charm ID", http.StatusBadRequest)
		return
	}

	// Fetch details of the weapon with the given ID
	eventURL := fmt.Sprintf("https://mhw-db.com/events/%s", eventID)
	var event Evenements
	err := fetchJSONData(eventURL, &event)
	if err != nil {
		http.Error(w, "Failed to fetch charm details", http.StatusInternalServerError)
		log.Println("Failed to fetch charm details:", err)
		return
	}

	// Parse the HTML template for displaying weapon details
	t, err := template.ParseFiles("templates/event_details.html")
	if err != nil {
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		log.Println("Failed to parse template:", err)
		return
	}

	// Render the template with the weapon details
	err = t.Execute(w, event)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		log.Println("Failed to render template:", err)
		return
	}
}

func decoHandler(w http.ResponseWriter, r *http.Request) {
	var deco []Decorations
	decoURL := "https://mhw-db.com/decorations"
	err := fetchJSONData(decoURL, &deco)
	if err != nil {
		log.Fatal(err)
	}
	// Obtenir les paramètres de filtrage de la requête
	rarityFilter := r.URL.Query().Get("rarityFilter")
	query := r.URL.Query().Get("q")

	// Filtrer les décorations en fonction des critères de filtrage
	var filteredDeco []Decorations
	for _, decoration := range deco {
		// Vérifiez si la décoration correspond aux critères de filtrage
		if (rarityFilter == "" || strconv.Itoa(decoration.Rarity) == rarityFilter) &&
			(query == "" || strings.Contains(strings.ToLower(decoration.Name), strings.ToLower(query))) {
			filteredDeco = append(filteredDeco, decoration)
		}
	}

	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1 // Si le numéro de page est invalide ou non défini, affichez la première page par défaut
	}

	// Calculez l'indice de début et de fin pour paginer les décorations
	startIndex := (page - 1) * 10
	endIndex := startIndex + 10
	if endIndex > len(filteredDeco) {
		endIndex = len(filteredDeco)
	}

	// Sélectionnez les décorations à afficher sur cette page
	paginatedDeco := filteredDeco[startIndex:endIndex]

	// Préparez les données nécessaires pour la pagination
	hasNextPage := endIndex < len(filteredDeco)
	hasPreviousPage := page > 1
	nextPage := page + 1
	previousPage := page - 1

	// Chargez le modèle HTML avec les données paginées et les informations de pagination
	data := struct {
		Decorations     []Decorations
		HasNextPage     bool
		HasPreviousPage bool
		NextPage        int
		PreviousPage    int
		RarityFilter    string
		Query           string
	}{
		Decorations:     paginatedDeco,
		HasNextPage:     hasNextPage,
		HasPreviousPage: hasPreviousPage,
		NextPage:        nextPage,
		PreviousPage:    previousPage,
		RarityFilter:    rarityFilter,
		Query:           query,
	}

	t, err := template.ParseFiles("templates/décorations.html")
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}

func decoDetailsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the weapon ID from the query parameters
	decoID := r.URL.Query().Get("id")
	if decoID == "" {
		http.Error(w, "Missing charm ID", http.StatusBadRequest)
		return
	}

	// Fetch details of the weapon with the given ID
	decoURL := fmt.Sprintf("https://mhw-db.com/decorations/%s", decoID)
	var deco Decorations
	err := fetchJSONData(decoURL, &deco)
	if err != nil {
		http.Error(w, "Failed to fetch charm details", http.StatusInternalServerError)
		log.Println("Failed to fetch charm details:", err)
		return
	}

	// Parse the HTML template for displaying weapon details
	t, err := template.ParseFiles("templates/deco_details.html")
	if err != nil {
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		log.Println("Failed to parse template:", err)
		return
	}

	// Render the template with the weapon details
	err = t.Execute(w, deco)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		log.Println("Failed to render template:", err)
		return
	}
}
func charmsHandler(w http.ResponseWriter, r *http.Request) {
	var charms []Charms
	charmsURL := "https://mhw-db.com/charms"
	err := fetchJSONData(charmsURL, &charms)
	if err != nil {
		log.Fatal(err)
	}

	// Obtenir le paramètre de recherche de la requête
	query := r.URL.Query().Get("q")

	// Filtrer les charms en fonction de la recherche
	var filteredCharms []Charms
	for _, charm := range charms {
		if query == "" || strings.Contains(strings.ToLower(charm.Name), strings.ToLower(query)) {
			filteredCharms = append(filteredCharms, charm)
		}
	}

	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1 // Si le numéro de page est invalide ou non défini, affichez la première page par défaut
	}

	// Calculez l'indice de début et de fin pour paginer les charms
	startIndex := (page - 1) * 10
	endIndex := startIndex + 10
	if endIndex > len(filteredCharms) {
		endIndex = len(filteredCharms)
	}

	// Sélectionnez les charms à afficher sur cette page
	paginatedCharms := filteredCharms[startIndex:endIndex]

	// Préparez les données nécessaires pour la pagination
	hasNextPage := endIndex < len(filteredCharms)
	hasPreviousPage := page > 1
	nextPage := page + 1
	previousPage := page - 1

	// Chargez le modèle HTML avec les données paginées et les informations de pagination
	data := struct {
		Charms          []Charms
		HasNextPage     bool
		HasPreviousPage bool
		NextPage        int
		PreviousPage    int
		Query           string
	}{
		Charms:          paginatedCharms,
		HasNextPage:     hasNextPage,
		HasPreviousPage: hasPreviousPage,
		NextPage:        nextPage,
		PreviousPage:    previousPage,
		Query:           query,
	}

	t, err := template.ParseFiles("templates/charms.html")
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}

func charmDetailsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the weapon ID from the query parameters
	charmID := r.URL.Query().Get("id")
	if charmID == "" {
		http.Error(w, "Missing charm ID", http.StatusBadRequest)
		return
	}

	// Fetch details of the weapon with the given ID
	charmURL := fmt.Sprintf("https://mhw-db.com/charms/%s", charmID)
	var charm Charms
	err := fetchJSONData(charmURL, &charm)
	if err != nil {
		http.Error(w, "Failed to fetch charm details", http.StatusInternalServerError)
		log.Println("Failed to fetch charm details:", err)
		return
	}

	// Parse the HTML template for displaying weapon details
	t, err := template.ParseFiles("templates/charm_details.html")
	if err != nil {
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		log.Println("Failed to parse template:", err)
		return
	}

	// Render the template with the weapon details
	err = t.Execute(w, charm)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		log.Println("Failed to render template:", err)
		return
	}
}

func weaponsHandler(w http.ResponseWriter, r *http.Request) {
	var weapons []Weapon
	weaponsURL := "https://mhw-db.com/weapons"
	err := fetchJSONData(weaponsURL, &weapons)
	if err != nil {
		log.Fatal(err)
	}

	// Obtenir les paramètres de filtrage de la requête
	rarityFilter := r.URL.Query().Get("rarityFilter")
	elementFilter := r.URL.Query().Get("elementFilter")
	damageTypeFilter := r.URL.Query().Get("damageTypeFilter")
	searchQuery := r.URL.Query().Get("q")

	// Filtrer les armes en fonction des critères de filtrage
	var filteredWeapons []Weapon
	for _, weapon := range weapons {
		// Vérifiez si l'arme correspond aux critères de filtrage
		if (rarityFilter == "" || strconv.Itoa(weapon.Rarity) == rarityFilter) &&
			(elementFilter == "" || containsElement(weapon.Element, elementFilter)) &&
			(damageTypeFilter == "" || weapon.TypeDegats == damageTypeFilter) &&
			(searchQuery == "" || strings.Contains(strings.ToLower(weapon.Name), strings.ToLower(searchQuery))) {
			filteredWeapons = append(filteredWeapons, weapon)
		}
	}

	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1 // Si le numéro de page est invalide ou non défini, affichez la première page par défaut
	}

	// Calculez l'indice de début et de fin pour paginer les armes
	startIndex := (page - 1) * 10
	endIndex := startIndex + 10
	if endIndex > len(filteredWeapons) {
		endIndex = len(filteredWeapons)
	}

	// Sélectionnez les armes à afficher sur cette page
	paginatedWeapons := filteredWeapons[startIndex:endIndex]

	// Préparez les données nécessaires pour la pagination
	hasNextPage := endIndex < len(filteredWeapons)
	hasPreviousPage := page > 1
	nextPage := page + 1
	previousPage := page - 1

	// Chargez le modèle HTML avec les données paginées et les informations de pagination
	data := struct {
		Weapons           []Weapon
		HasNextPage       bool
		HasPreviousPage   bool
		NextPage          int
		PreviousPage      int
		Query             string
		RarityFilter      string
		ElementFilter     string
		DamagesTypeFilter string
	}{
		Weapons:           paginatedWeapons,
		HasNextPage:       hasNextPage,
		HasPreviousPage:   hasPreviousPage,
		NextPage:          nextPage,
		PreviousPage:      previousPage,
		Query:             searchQuery,
		RarityFilter:      rarityFilter,
		ElementFilter:     elementFilter,
		DamagesTypeFilter: damageTypeFilter,
	}

	t, err := template.ParseFiles("templates/weapons.html")
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}

func containsElement(elements []struct {
	Type   string "json:\"type\""
	Damage int    "json:\"damage\""
}, elementFilter string) bool {
	for _, element := range elements {
		if element.Type == elementFilter {
			return true
		}
	}
	return false
}

func weaponDetailsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the weapon ID from the query parameters
	weaponID := r.URL.Query().Get("id")
	if weaponID == "" {
		http.Error(w, "Missing weapon ID", http.StatusBadRequest)
		return
	}

	// Fetch details of the weapon with the given ID
	weaponURL := fmt.Sprintf("https://mhw-db.com/weapons/%s", weaponID)
	var weapon Weapon
	err := fetchJSONData(weaponURL, &weapon)
	if err != nil {
		http.Error(w, "Failed to fetch weapon details", http.StatusInternalServerError)
		log.Println("Failed to fetch weapon details:", err)
		return
	}

	// Parse the HTML template for displaying weapon details
	t, err := template.ParseFiles("templates/weapon_details.html")
	if err != nil {
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		log.Println("Failed to parse template:", err)
		return
	}

	// Render the template with the weapon details
	err = t.Execute(w, weapon)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		log.Println("Failed to render template:", err)
		return
	}
}

func armorsHandler(w http.ResponseWriter, r *http.Request) {

	armorsURL := "https://mhw-db.com/armor"

	err := fetchJSONData(armorsURL, &armors)
	if err != nil {
		log.Fatal(err)
	}

	// Obtenir les paramètres de filtrage de la requête
	typeFilter := r.URL.Query().Get("typeFilter")
	rankFilter := r.URL.Query().Get("rankFilter")
	rarityFilter := r.URL.Query().Get("rarityFilter")
	query := r.URL.Query().Get("q")

	// Filtrer les armures en fonction des critères de filtrage
	var filteredArmors []Armor
	for _, armor := range armors {
		// Vérifiez si l'armure correspond aux critères de filtrage
		if (typeFilter == "" || armor.Type == typeFilter) &&
			(rankFilter == "" || armor.Rang == rankFilter) &&
			(rarityFilter == "" || strconv.Itoa(armor.Rarity) == rarityFilter) &&
			(query == "" || strings.Contains(strings.ToLower(armor.Name), strings.ToLower(query))) {
			filteredArmors = append(filteredArmors, armor)
		}
	}

	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1 // Si le numéro de page est invalide ou non défini, affichez la première page par défaut
	}

	// Calculez l'indice de début et de fin pour paginer les armures
	startIndex := (page - 1) * 10
	endIndex := startIndex + 10
	if endIndex > len(filteredArmors) {
		endIndex = len(filteredArmors)
	}

	// Sélectionnez les armures à afficher sur cette page
	paginatedArmors := filteredArmors[startIndex:endIndex]

	// Préparez les données nécessaires pour la pagination
	hasNextPage := endIndex < len(filteredArmors)
	hasPreviousPage := page > 1
	nextPage := page + 1
	previousPage := page - 1

	// Chargez le modèle HTML avec les données paginées et les informations de pagination
	data := struct {
		Armors          []Armor
		HasNextPage     bool
		HasPreviousPage bool
		NextPage        int
		PreviousPage    int
		TypeFilter      string
		RankFilter      string
		RarityFilter    string
		Query           string
	}{
		Armors:          paginatedArmors,
		HasNextPage:     hasNextPage,
		HasPreviousPage: hasPreviousPage,
		NextPage:        nextPage,
		PreviousPage:    previousPage,
		TypeFilter:      typeFilter,
		RankFilter:      rankFilter,
		RarityFilter:    rarityFilter,
		Query:           query,
	}

	t, err := template.ParseFiles("templates/armors.html")
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}

func armorDetailsHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID de l'armure à partir des paramètres de requête
	armorID := r.URL.Query().Get("id")
	fmt.Println("Armor URL Path:", r.URL.Path)            // Point de contrôle
	fmt.Println("Trying to find armor with ID:", armorID) // Ajout du point de contrôle

	// Convertir l'ID en entier
	id, err := strconv.Atoi(armorID)
	if err != nil {
		http.Error(w, "Invalid armor ID", http.StatusBadRequest)
		return
	}

	// Récupérer les détails de l'armure à partir de l'ID
	var armor Armor
	for _, a := range armors {
		if a.ID == id {
			armor = a
			break
		}
	}

	fmt.Println("Found armor with ID:", armor.ID)
	// Vérifier si l'armure a été trouvée
	if armor.ID == 0 {
		http.Error(w, "Armor not found", http.StatusNotFound)
		return
	}

	// Charger le template et afficher les détails de l'armure
	t, err := template.ParseFiles("templates/armor_details.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, armor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func monstersHandler(w http.ResponseWriter, r *http.Request) {
	var monsters []Monster
	monstersURL := "https://mhw-db.com/monsters"
	err := fetchJSONData(monstersURL, &monsters)
	if err != nil {
		log.Fatal(err)
	}

	// Obtenir les paramètres de filtrage de la requête
	sizeFilter := r.URL.Query().Get("sizeFilter")
	resistanceFilter := r.URL.Query().Get("resistanceFilter")
	searchQuery := r.URL.Query().Get("q")

	// Filtrer les monstres en fonction des critères de filtrage
	var filteredMonsters []Monster
	for _, monster := range monsters {
		// Vérifiez si le monstre correspond aux critères de filtrage
		if (sizeFilter == "" || monster.Size == sizeFilter) &&
			(resistanceFilter == "" || containsResistance(monster.Resistances, resistanceFilter)) &&
			(searchQuery == "" || strings.Contains(strings.ToLower(monster.Name), strings.ToLower(searchQuery))) {
			filteredMonsters = append(filteredMonsters, monster)
		}
	}

	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1 // Si le numéro de page est invalide ou non défini, affichez la première page par défaut
	}

	// Calculez l'indice de début et de fin pour paginer les monstres
	startIndex := (page - 1) * 10
	endIndex := startIndex + 10
	if endIndex > len(filteredMonsters) {
		endIndex = len(filteredMonsters)
	}

	// Sélectionnez les monstres à afficher sur cette page
	paginatedMonsters := filteredMonsters[startIndex:endIndex]

	// Préparez les données nécessaires pour la pagination
	hasNextPage := endIndex < len(filteredMonsters)
	hasPreviousPage := page > 1
	nextPage := page + 1
	previousPage := page - 1

	// Chargez le modèle HTML avec les données paginées et les informations de pagination
	data := struct {
		Monsters         []Monster
		HasNextPage      bool
		HasPreviousPage  bool
		NextPage         int
		PreviousPage     int
		Query            string
		SizeFilter       string
		ResistanceFilter string
	}{
		Monsters:         paginatedMonsters,
		HasNextPage:      hasNextPage,
		HasPreviousPage:  hasPreviousPage,
		NextPage:         nextPage,
		PreviousPage:     previousPage,
		Query:            searchQuery,
		SizeFilter:       sizeFilter,
		ResistanceFilter: resistanceFilter,
	}

	t, err := template.ParseFiles("templates/monsters.html")
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}

func containsResistance(resistances []struct {
	Element string "json:\"element\""
}, resistanceFilter string) bool {
	for _, resistance := range resistances {
		if resistance.Element == resistanceFilter {
			return true
		}
	}
	return false
}

func monsterDetailsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the weapon ID from the query parameters
	monsterID := r.URL.Query().Get("id")
	if monsterID == "" {
		http.Error(w, "Missing weapon ID", http.StatusBadRequest)
		return
	}

	// Fetch details of the weapon with the given ID
	monsterURL := fmt.Sprintf("https://mhw-db.com/monsters/%s", monsterID)
	var monster Monster
	err := fetchJSONData(monsterURL, &monster)
	if err != nil {
		http.Error(w, "Failed to fetch weapon details", http.StatusInternalServerError)
		log.Println("Failed to fetch weapon details:", err)
		return
	}

	// Parse the HTML template for displaying weapon details
	t, err := template.ParseFiles("templates/monster_details.html")
	if err != nil {
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		log.Println("Failed to parse template:", err)
		return
	}

	// Render the template with the weapon details
	err = t.Execute(w, monster)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		log.Println("Failed to render template:", err)
		return
	}
}

// Fonction pour supprimer des éléments des favoris
func removeFromFavoritesHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer les paramètres de l'URL pour identifier l'élément à supprimer
	itemID := r.URL.Query().Get("id")
	itemType := r.URL.Query().Get("type")

	// Vérifier le type d'élément et agir en conséquence
	switch itemType {
	case "weapons":
		// Parcourir la liste des armes favorites pour trouver l'élément à supprimer
		for i, weapon := range favoriteWeapons {
			if strconv.Itoa(weapon.ID) == itemID {
				// Supprimer l'élément de la liste des armes favorites
				favoriteWeapons = append(favoriteWeapons[:i], favoriteWeapons[i+1:]...)
				break
			}
		}
	case "monsters":
		// Parcourir la liste des monstres favoris pour trouver l'élément à supprimer
		for i, monster := range favoriteMonsters {
			if strconv.Itoa(monster.ID) == itemID {
				// Supprimer l'élément de la liste des monstres favoris
				favoriteMonsters = append(favoriteMonsters[:i], favoriteMonsters[i+1:]...)
				break
			}
		}
	case "items":
		// Parcourir la liste des objets favoris pour trouver l'élément à supprimer
		for i, item := range favoriteItems {
			if strconv.Itoa(item.ID) == itemID && item.Category == "items" {
				// Supprimer l'élément de la liste des objets favoris
				favoriteItems = append(favoriteItems[:i], favoriteItems[i+1:]...)
				break
			}
		}
	case "skills":
		// Parcourir la liste des compétences favorites pour trouver l'élément à supprimer
		for i, skill := range favoriteSkills {
			if strconv.Itoa(skill.ID) == itemID {
				// Supprimer l'élément de la liste des compétences favorites
				favoriteSkills = append(favoriteSkills[:i], favoriteSkills[i+1:]...)
				break
			}
		}
	case "events":
		// Parcourir la liste des événements favoris pour trouver l'élément à supprimer
		for i, event := range favoriteEvents {
			if strconv.Itoa(event.ID) == itemID {
				// Supprimer l'élément de la liste des événements favoris
				favoriteEvents = append(favoriteEvents[:i], favoriteEvents[i+1:]...)
				break
			}
		}
	case "deco":
		// Parcourir la liste des décorations favorites pour trouver l'élément à supprimer
		for i, deco := range favoriteDeco {
			if strconv.Itoa(deco.ID) == itemID {
				// Supprimer l'élément de la liste des décorations favorites
				favoriteDeco = append(favoriteDeco[:i], favoriteDeco[i+1:]...)
				break
			}
		}
	case "charms":
		// Parcourir la liste des charmes favoris pour trouver l'élément à supprimer
		for i, charm := range favoriteCharms {
			if strconv.Itoa(charm.ID) == itemID {
				// Supprimer l'élément de la liste des charmes favoris
				favoriteCharms = append(favoriteCharms[:i], favoriteCharms[i+1:]...)
				break
			}
		}
	case "armors":
		// Parcourir la liste des armures favorites pour trouver l'élément à supprimer
		for i, armor := range favoriteArmors {
			if strconv.Itoa(armor.ID) == itemID {
				// Supprimer l'élément de la liste des armures favorites
				favoriteArmors = append(favoriteArmors[:i], favoriteArmors[i+1:]...)
				break
			}
		}
	default:
		// Type d'élément non pris en charge
		http.Error(w, "Unsupported item type", http.StatusBadRequest)
		return
	}

	// Rediriger l'utilisateur vers la page précédente (page de favoris)
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
}

func addToFavoritesHandler(w http.ResponseWriter, r *http.Request) {
	itemCate := r.URL.Query().Get("category")
	itemID := r.URL.Query().Get("id")
	itemName := r.URL.Query().Get("name")
	itemType := r.URL.Query().Get("type")
	itemRarity := r.URL.Query().Get("rarity")
	itemImage := r.URL.Query().Get("image")
	itemImageH := r.URL.Query().Get("imageMale")
	itemImageF := r.URL.Query().Get("imageFemale")
	itemSpecies := r.URL.Query().Get("species")

	// Vérifiez le type d'élément et agissez en conséquence
	switch itemCate {
	case "weapons":
		// Ajoutez l'arme aux favoris
		var weapon Weapon
		weapon.ID, _ = strconv.Atoi(itemID)
		weapon.Name = itemName
		weapon.Type = itemType
		weapon.Rarity, _ = strconv.Atoi(itemRarity)
		weapon.Assets.Image = itemImage
		// Ajoutez d'autres détails de l'arme si nécessaire
		// ...

		// Ajoutez l'arme aux favoris
		weapons := Weapon{
			Category: "weapons",
			ID:       weapon.ID,
			Name:     weapon.Name,
			Type:     weapon.Type,
			Rarity:   weapon.Rarity,
			Assets: struct {
				Image string `json:"image"`
			}{
				Image: itemImage,
			},
			// Ajoutez d'autres détails de l'arme si nécessaire
			// ...
		}
		// Ajoutez l'arme aux favoris
		favoriteWeapons = append(favoriteWeapons, weapons)
	case "monsters":
		// Ajoutez le monstre aux favoris
		var monster Monster
		monster.ID, _ = strconv.Atoi(itemID)
		monster.Name = itemName
		monster.Species = itemSpecies
		monster.Size = itemType
		// Ajoutez d'autres détails du monstre si nécessaire
		// ...

		// Ajoutez le monstre aux favoris
		monsters := Monster{
			Category: "monsters",
			ID:       monster.ID,
			Name:     monster.Name,
			Size:     monster.Size,
			Species:  monster.Species,
			// Ajoutez d'autres détails du monstre si nécessaire
			// ...
		}
		// Ajoutez le monstre aux favoris
		favoriteMonsters = append(favoriteMonsters, monsters)
	case "items":
		// Ajoutez le monstre aux favoris
		var objet Objets
		objet.ID, _ = strconv.Atoi(itemID)
		objet.Name = itemName
		objet.Rarity, _ = strconv.Atoi(itemRarity)

		// Ajoutez d'autres détails du monstre si nécessaire
		// ...

		// Ajoutez le monstre aux favoris
		objets := Objets{
			Category: "items",
			ID:       objet.ID,
			Name:     objet.Name,
			Rarity:   objet.Rarity,
			// Ajoutez d'autres détails du monstre si nécessaire
			// ...
		}
		// Ajoutez le monstre aux favoris
		favoriteObjets = append(favoriteObjets, objets)
	case "skills":
		var skill Skills
		skill.ID, _ = strconv.Atoi(itemID)
		skill.Name = itemName

		// Ajoutez d'autres détails de l'arme si nécessaire
		// ...

		// Ajoutez l'arme aux favoris
		skills := Skills{
			Category: "skills",
			ID:       skill.ID,
			Name:     skill.Name,

			// Ajoutez d'autres détails de l'arme si nécessaire
			// ...
		}
		// Ajoutez l'arme aux favoris
		favoriteSkills = append(favoriteSkills, skills)
	case "events":
		var event Evenements
		event.ID, _ = strconv.Atoi(itemID)
		event.Name = itemName
		event.Type = itemType

		// Ajoutez d'autres détails de l'arme si nécessaire
		// ...

		// Ajoutez l'arme aux favoris
		events := Evenements{
			Category: "events",
			ID:       event.ID,
			Name:     event.Name,
			Type:     event.Type,

			// Ajoutez d'autres détails de l'arme si nécessaire
			// ...
		}
		// Ajoutez l'arme aux favoris
		favoriteEvents = append(favoriteEvents, events)
	case "deco":
		var deco Decorations
		deco.ID, _ = strconv.Atoi(itemID)
		deco.Name = itemName
		deco.Rarity, _ = strconv.Atoi(itemRarity)

		// Ajoutez d'autres détails de l'arme si nécessaire
		// ...

		// Ajoutez l'arme aux favoris
		decos := Decorations{
			Category: "deco",
			ID:       deco.ID,
			Name:     deco.Name,
			Rarity:   deco.Rarity,

			// Ajoutez d'autres détails de l'arme si nécessaire
			// ...
		}
		// Ajoutez l'arme aux favoris
		favoriteDeco = append(favoriteDeco, decos)
	case "charms":
		var charm Charms
		charm.ID, _ = strconv.Atoi(itemID)
		charm.Name = itemName

		// Ajoutez d'autres détails de l'arme si nécessaire
		// ...

		// Ajoutez l'arme aux favoris
		charms := Charms{
			Category: "charms",
			ID:       charm.ID,
			Name:     charm.Name,

			// Ajoutez d'autres détails de l'arme si nécessaire
			// ...
		}
		// Ajoutez l'arme aux favoris
		favoriteCharms = append(favoriteCharms, charms)
	case "armors":
		// Ajoutez l'arme aux favoris
		var armor Armor
		armor.ID, _ = strconv.Atoi(itemID)
		armor.Name = itemName
		armor.Type = itemType
		armor.Rarity, _ = strconv.Atoi(itemRarity)
		armor.Assets.ImageH = itemImageH
		armor.Assets.ImageF = itemImageF
		// Ajoutez d'autres détails de l'arme si nécessaire
		// ...

		// Ajoutez l'arme aux favoris
		armors := Armor{
			Category: "armors",
			ID:       armor.ID,
			Name:     armor.Name,
			Type:     armor.Type,
			Rarity:   armor.Rarity,
			Assets: struct {
				ImageH string `json:"imageMale"`
				ImageF string `json:"imageFemale"`
			}{
				ImageH: itemImageH,
				ImageF: itemImageF,
			},
			// Ajoutez d'autres détails de l'arme si nécessaire
			// ...
		}
		// Ajoutez l'arme aux favoris
		favoriteArmors = append(favoriteArmors, armors)

	default:
		// Type d'élément non pris en charge
		http.Error(w, "Unsupported item type", http.StatusBadRequest)
		return
	}

	// Redirigez l'utilisateur vers la page précédente
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
}

func favoritesHandler(w http.ResponseWriter, r *http.Request) {
	// Chargez le modèle HTML pour afficher les favoris
	t, err := template.ParseFiles("templates/favorites.html")
	if err != nil {
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		log.Println("Failed to parse template:", err)
		return
	}
	fmt.Println("Favorite items:", favoriteItems)
	// Préparez les données des favoris pour le rendu dans le modèle
	data := struct {
		FavoriteItems Favorites
	}{
		FavoriteItems: Favorites{
			Weapons:  favoriteWeapons,
			Monsters: favoriteMonsters,
			Skills:   favoriteSkills,
			Objets:   favoriteObjets,
			Events:   favoriteEvents,
			Deco:     favoriteDeco,
			Charms:   favoriteCharms,
			Armors:   favoriteArmors,
		},
	}

	// Rendez le modèle HTML avec les données des favoris
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		log.Println("Failed to render template:", err)
		return
	}
}

func main() {
	// Serveur HTTP
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", indexHandler)

	http.HandleFunc("/search", mainPageHandler)

	http.HandleFunc("/skills", skillsHandler)
	http.HandleFunc("/skill-details", skillDetailsHandler)

	http.HandleFunc("/items", itemsHandler)
	http.HandleFunc("/item-details", itemDetailsHandler)

	http.HandleFunc("/events", eventsHandler)
	http.HandleFunc("/event-details", eventDetailsHandler)

	http.HandleFunc("/décorations", decoHandler)
	http.HandleFunc("/deco-details", decoDetailsHandler)

	http.HandleFunc("/charms", charmsHandler)
	http.HandleFunc("/charm-details", charmDetailsHandler)

	http.HandleFunc("/armors", armorsHandler)
	http.HandleFunc("/armor-details", armorDetailsHandler) // Nouvelle route pour afficher les détails de l'armure

	http.HandleFunc("/weapons", weaponsHandler)
	http.HandleFunc("/weapon-details", weaponDetailsHandler)

	http.HandleFunc("/monsters", monstersHandler)
	http.HandleFunc("/monster-details", monsterDetailsHandler)

	http.HandleFunc("/favorites", favoritesHandler)
	http.HandleFunc("/add-favorite", addToFavoritesHandler)
	http.HandleFunc("/remove-favorite", removeFromFavoritesHandler)

	log.Println("Serveur en écoute sur le port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
