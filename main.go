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

// Structure pour stocker les données des armes
type Weapon struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Rarity int    `json:"rarity"`
	Attack struct {
		Display int `json:"display"`
		Raw     int `json:"raw"`
	} `json:"attack"`
	Element []struct {
		Type   string `json:"type"`
		Damage int    `json:"damage"`
	} `json:"elements"`
	TypeDegats string `json:"damageType"`
	Assets     struct {
		Image string `json:"image"`
	} `json:"assets"`
}

type Skills struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Rang        []struct {
		Niveau      int    `json:"level"`
		Description string `json:"description"`
	} `json:"ranks"`
}

type Objets struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Rarete      int    `json:"rarity"`
	Description string `json:"description"`
}

type Decorations struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Rarete int    `json:"rarity"`
	Skills []struct {
		Name        string `json:"skillName"`
		Niveau      int    `json:"level"`
		Description string `json:"description"`
	} `json:"skills"`
}

type Charms struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Rang []struct {
		Niveau int `json:"level"`
		Rarete int `json:"rarity"`
		Skills []struct {
			Name        string `json:"skillName"`
			Niveau      int    `json:"level"`
			Description string `json:"description"`
		}
	} `json:"ranks"`
}

type Armor struct {
	ID          int            `json:"id"`
	Name        string         `json:"name"`
	Type        string         `json:"type"`
	Rarete      int            `json:"rarity"`
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
}

// Structure pour stocker les données des monstres
type Monster struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Species     string `json:"species"`
	Type        string `json:"type"`
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
}

type Evenements struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Rang         int    `json:"questRank"`
	Description  string `json:"description"`
	Exigences    string `json:"requirements"`
	Objectif     string `json:"successConditions"`
	Localisation struct {
		Name string `json:"name"`
	} `json:"location"`
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
	// Filtrer les compétences en fonction de la recherche
	query := r.URL.Query().Get("q")
	if query != "" {
		var filteredSkills []Skills
		for _, skill := range skills {
			if strings.Contains(strings.ToLower(skill.Name), strings.ToLower(query)) {
				filteredSkills = append(filteredSkills, skill)
			}
		}
		skills = filteredSkills
	}
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1 // Si le numéro de page est invalide ou non défini, affichez la première page par défaut
	}

	// Calculez l'indice de début et de fin pour paginer les armes
	startIndex := (page - 1) * 10
	endIndex := startIndex + 10
	if endIndex > len(skills) {
		endIndex = len(skills)
	}

	// Sélectionnez les armes à afficher sur cette page
	paginatedSkills := skills[startIndex:endIndex]

	// Préparez les données nécessaires pour la pagination
	hasNextPage := endIndex < len(skills)
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
	}{
		Skills:          paginatedSkills,
		HasNextPage:     hasNextPage,
		HasPreviousPage: hasPreviousPage,
		NextPage:        nextPage,
		PreviousPage:    previousPage,
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
	// Filtrer les objets en fonction de la recherche
	query := r.URL.Query().Get("q")
	if query != "" {
		var filteredItems []Objets
		for _, item := range items {
			if strings.Contains(strings.ToLower(item.Name), strings.ToLower(query)) {
				filteredItems = append(filteredItems, item)
			}
		}
		items = filteredItems
	}
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1 // Si le numéro de page est invalide ou non défini, affichez la première page par défaut
	}

	// Calculez l'indice de début et de fin pour paginer les armes
	startIndex := (page - 1) * 10
	endIndex := startIndex + 10
	if endIndex > len(items) {
		endIndex = len(items)
	}

	// Sélectionnez les armes à afficher sur cette page
	paginatedItems := items[startIndex:endIndex]

	// Préparez les données nécessaires pour la pagination
	hasNextPage := endIndex < len(items)
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
	}{
		Items:           paginatedItems,
		HasNextPage:     hasNextPage,
		HasPreviousPage: hasPreviousPage,
		NextPage:        nextPage,
		PreviousPage:    previousPage,
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
	// Filtrer les événements en fonction de la recherche
	query := r.URL.Query().Get("q")
	if query != "" {
		var filteredEvents []Evenements
		for _, event := range events {
			if strings.Contains(strings.ToLower(event.Name), strings.ToLower(query)) {
				filteredEvents = append(filteredEvents, event)
			}
		}
		events = filteredEvents
	}
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1 // Si le numéro de page est invalide ou non défini, affichez la première page par défaut
	}

	// Calculez l'indice de début et de fin pour paginer les armes
	startIndex := (page - 1) * 10
	endIndex := startIndex + 10
	if endIndex > len(events) {
		endIndex = len(events)
	}

	// Sélectionnez les armes à afficher sur cette page
	paginatedEvents := events[startIndex:endIndex]

	// Préparez les données nécessaires pour la pagination
	hasNextPage := endIndex < len(events)
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
	}{
		Events:          paginatedEvents,
		HasNextPage:     hasNextPage,
		HasPreviousPage: hasPreviousPage,
		NextPage:        nextPage,
		PreviousPage:    previousPage,
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
	// Filtrer les décorations en fonction de la recherche
	query := r.URL.Query().Get("q")
	if query != "" {
		var filteredDeco []Decorations
		for _, decoration := range deco {
			if strings.Contains(strings.ToLower(decoration.Name), strings.ToLower(query)) {
				filteredDeco = append(filteredDeco, decoration)
			}
		}
		deco = filteredDeco
	}
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1 // Si le numéro de page est invalide ou non défini, affichez la première page par défaut
	}

	// Calculez l'indice de début et de fin pour paginer les armes
	startIndex := (page - 1) * 10
	endIndex := startIndex + 10
	if endIndex > len(deco) {
		endIndex = len(deco)
	}

	// Sélectionnez les armes à afficher sur cette page
	paginatedDeco := deco[startIndex:endIndex]

	// Préparez les données nécessaires pour la pagination
	hasNextPage := endIndex < len(deco)
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
	}{
		Decorations:     paginatedDeco,
		HasNextPage:     hasNextPage,
		HasPreviousPage: hasPreviousPage,
		NextPage:        nextPage,
		PreviousPage:    previousPage,
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
	// Filtrer les charms en fonction de la recherche
	query := r.URL.Query().Get("q")
	if query != "" {
		var filteredCharms []Charms
		for _, charm := range charms {
			if strings.Contains(strings.ToLower(charm.Name), strings.ToLower(query)) {
				filteredCharms = append(filteredCharms, charm)
			}
		}
		charms = filteredCharms
	}
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1 // Si le numéro de page est invalide ou non défini, affichez la première page par défaut
	}

	// Calculez l'indice de début et de fin pour paginer les armes
	startIndex := (page - 1) * 10
	endIndex := startIndex + 10
	if endIndex > len(charms) {
		endIndex = len(charms)
	}

	// Sélectionnez les armes à afficher sur cette page
	paginatedCharms := charms[startIndex:endIndex]

	// Préparez les données nécessaires pour la pagination
	hasNextPage := endIndex < len(charms)
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
	}{
		Charms:          paginatedCharms,
		HasNextPage:     hasNextPage,
		HasPreviousPage: hasPreviousPage,
		NextPage:        nextPage,
		PreviousPage:    previousPage,
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
		Weapons         []Weapon
		HasNextPage     bool
		HasPreviousPage bool
		NextPage        int
		PreviousPage    int
		Query           string
	}{
		Weapons:         paginatedWeapons,
		HasNextPage:     hasNextPage,
		HasPreviousPage: hasPreviousPage,
		NextPage:        nextPage,
		PreviousPage:    previousPage,
		Query:           searchQuery,
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

	// Filtrer les armures en fonction de la recherche
	query := r.URL.Query().Get("q")
	if query != "" {
		var filteredArmors []Armor
		for _, armor := range armors {
			if strings.Contains(strings.ToLower(armor.Name), strings.ToLower(query)) {
				filteredArmors = append(filteredArmors, armor)
			}
		}
		armors = filteredArmors
	}
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1 // Si le numéro de page est invalide ou non défini, affichez la première page par défaut
	}

	// Calculez l'indice de début et de fin pour paginer les armes
	startIndex := (page - 1) * 10
	endIndex := startIndex + 10
	if endIndex > len(armors) {
		endIndex = len(armors)
	}

	// Sélectionnez les armes à afficher sur cette page
	paginatedArmors := armors[startIndex:endIndex]

	// Préparez les données nécessaires pour la pagination
	hasNextPage := endIndex < len(armors)
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
	}{
		Armors:          paginatedArmors,
		HasNextPage:     hasNextPage,
		HasPreviousPage: hasPreviousPage,
		NextPage:        nextPage,
		PreviousPage:    previousPage,
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

	// Filtrer les monstres en fonction de la taille
	sizeFilter := r.URL.Query().Get("sizeFilter")
	if sizeFilter != "" {
		var filteredMonsters []Monster
		for _, monster := range monsters {
			if monster.Size == sizeFilter {
				filteredMonsters = append(filteredMonsters, monster)
			}
		}
		monsters = filteredMonsters
	}

	// Filtrer les monstres en fonction des résistances
	resistanceFilter := r.URL.Query().Get("resistanceFilter")
	if resistanceFilter != "" {
		var filteredMonsters []Monster
		for _, monster := range monsters {
			if containsResistance(monster.Resistances, resistanceFilter) {
				filteredMonsters = append(filteredMonsters, monster)
			}
		}
		monsters = filteredMonsters
	}

	// Filtrer les monstres en fonction de la recherche
	query := r.URL.Query().Get("q")
	if query != "" {
		var filteredMonsters []Monster
		for _, monster := range monsters {
			if strings.Contains(strings.ToLower(monster.Name), strings.ToLower(query)) {
				filteredMonsters = append(filteredMonsters, monster)
			}
		}
		monsters = filteredMonsters
	}

	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1 // Si le numéro de page est invalide ou non défini, affichez la première page par défaut
	}

	// Calculez l'indice de début et de fin pour paginer les monstres
	startIndex := (page - 1) * 10
	endIndex := startIndex + 10
	if endIndex > len(monsters) {
		endIndex = len(monsters)
	}

	// Sélectionnez les monstres à afficher sur cette page
	paginatedMonsters := monsters[startIndex:endIndex]

	// Préparez les données nécessaires pour la pagination
	hasNextPage := endIndex < len(monsters)
	hasPreviousPage := page > 1
	nextPage := page + 1
	previousPage := page - 1

	// Chargez le modèle HTML avec les données paginées et les informations de pagination
	data := struct {
		Monsters        []Monster
		HasNextPage     bool
		HasPreviousPage bool
		NextPage        int
		PreviousPage    int
	}{
		Monsters:        paginatedMonsters,
		HasNextPage:     hasNextPage,
		HasPreviousPage: hasPreviousPage,
		NextPage:        nextPage,
		PreviousPage:    previousPage,
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

// Helper function to check if a monster has a specific resistance
func containsResistance(resistances []Resistance, resistanceFilter string) bool {
	for _, resistance := range resistances {
		if resistance.Element == resistanceFilter {
			return true
		}
	}
	return false
}

// func monstersHandler(w http.ResponseWriter, r *http.Request) {
// 	var monsters []Monster
// 	monstersURL := "https://mhw-db.com/monsters"
// 	err := fetchJSONData(monstersURL, &monsters)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	// Filtrer les monstres en fonction de la recherche
// 	query := r.URL.Query().Get("q")
// 	if query != "" {
// 		var filteredMonsters []Monster
// 		for _, monster := range monsters {
// 			if strings.Contains(strings.ToLower(monster.Name), strings.ToLower(query)) {
// 				filteredMonsters = append(filteredMonsters, monster)
// 			}
// 		}
// 		monsters = filteredMonsters
// 	}

// 	pageStr := r.URL.Query().Get("page")
// 	page, err := strconv.Atoi(pageStr)
// 	if err != nil || page < 1 {
// 		page = 1 // Si le numéro de page est invalide ou non défini, affichez la première page par défaut
// 	}

// 	// Calculez l'indice de début et de fin pour paginer les armes
// 	startIndex := (page - 1) * 10
// 	endIndex := startIndex + 10
// 	if endIndex > len(monsters) {
// 		endIndex = len(monsters)
// 	}

// 	// Sélectionnez les armes à afficher sur cette page
// 	paginatedMonsters := monsters[startIndex:endIndex]

// 	// Préparez les données nécessaires pour la pagination
// 	hasNextPage := endIndex < len(monsters)
// 	hasPreviousPage := page > 1
// 	nextPage := page + 1
// 	previousPage := page - 1

// 	// Chargez le modèle HTML avec les données paginées et les informations de pagination
// 	data := struct {
// 		Monsters        []Monster
// 		HasNextPage     bool
// 		HasPreviousPage bool
// 		NextPage        int
// 		PreviousPage    int
// 	}{
// 		Monsters:        paginatedMonsters,
// 		HasNextPage:     hasNextPage,
// 		HasPreviousPage: hasPreviousPage,
// 		NextPage:        nextPage,
// 		PreviousPage:    previousPage,
// 	}

// 	t, err := template.ParseFiles("templates/monsters.html")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	err = t.Execute(w, data)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

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
func main() {
	// Serveur HTTP
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", indexHandler)

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

	log.Println("Serveur en écoute sur le port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
