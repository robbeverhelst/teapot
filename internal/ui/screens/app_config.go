package screens

import (
	"teapot/internal/models"
	"teapot/internal/ui/components"
	"teapot/internal/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AppConfigOption struct {
	Key         string
	Name        string
	Description string
	Selected    bool
}

type AppConfigModel struct {
	appType     models.AppType
	options     []AppConfigOption
	cursor      int
	appName     string
	nameCursor  int
	focusedArea int // 0 = name, 1 = options
}

func NewAppConfigModel(appType models.AppType) AppConfigModel {
	options := getOptionsForAppType(appType)
	
	return AppConfigModel{
		appType:     appType,
		options:     options,
		cursor:      0,
		appName:     getDefaultAppName(appType),
		nameCursor:  len(getDefaultAppName(appType)),
		focusedArea: 0,
	}
}

func (m AppConfigModel) Init() tea.Cmd {
	return nil
}

func (m AppConfigModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			m.focusedArea = (m.focusedArea + 1) % 2
		case "shift+tab":
			m.focusedArea = (m.focusedArea - 1 + 2) % 2
		case "enter":
			if m.appName != "" {
				// Collect selected options
				selectedOptions := make(map[string]interface{})
				for _, option := range m.options {
					selectedOptions[option.Key] = option.Selected
				}
				
				return m, func() tea.Msg {
					return AppConfigCompleteMsg{
						AppName: m.appName,
						Options: selectedOptions,
					}
				}
			}
		case "j", "down":
			if m.focusedArea == 1 && m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case "k", "up":
			if m.focusedArea == 1 && m.cursor > 0 {
				m.cursor--
			}
		case " ":
			if m.focusedArea == 1 {
				m.options[m.cursor].Selected = !m.options[m.cursor].Selected
			}
		case "backspace":
			if m.focusedArea == 0 && m.nameCursor > 0 {
				m.appName = m.appName[:m.nameCursor-1] + m.appName[m.nameCursor:]
				m.nameCursor--
			}
		case "left":
			if m.focusedArea == 0 && m.nameCursor > 0 {
				m.nameCursor--
			}
		case "right":
			if m.focusedArea == 0 && m.nameCursor < len(m.appName) {
				m.nameCursor++
			}
		case "home":
			if m.focusedArea == 0 {
				m.nameCursor = 0
			}
		case "end":
			if m.focusedArea == 0 {
				m.nameCursor = len(m.appName)
			}
		default:
			if m.focusedArea == 0 && len(msg.String()) == 1 {
				char := msg.String()
				if isValidAppNameChar(char) {
					m.appName = m.appName[:m.nameCursor] + char + m.appName[m.nameCursor:]
					m.nameCursor++
				}
			}
		}
	}
	return m, nil
}

func (m AppConfigModel) View() string {
	appTypeName := models.AppTypeNames[m.appType]
	subtitle := components.RenderSubtitle("Configure your " + appTypeName + " application")

	// App name input
	nameLabel := lipgloss.NewStyle().
		Foreground(styles.ColorTextPrimary).
		Bold(true).
		Margin(1, 0, 0, 0).
		Render("📝 Application name:")

	nameValue := m.appName
	if m.focusedArea == 0 {
		if m.nameCursor < len(nameValue) {
			nameValue = nameValue[:m.nameCursor] + "_" + nameValue[m.nameCursor+1:]
		} else {
			nameValue += "_"
		}
	}

	nameBorderColor := styles.ColorBorderPrimary
	if m.focusedArea == 0 {
		nameBorderColor = styles.ColorBorderAccent
	}

	nameBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(nameBorderColor).
		Foreground(styles.ColorTextPrimary).
		Padding(0, 1).
		Width(35).
		Margin(1, 0, 0, 0).
		Render(nameValue)

	// Options section
	optionsLabel := lipgloss.NewStyle().
		Foreground(styles.ColorTextPrimary).
		Bold(true).
		Margin(2, 0, 0, 0).
		Render("⚙️ Features & Integrations:")

	var optionsList string
	for i, option := range m.options {
		cursor := " "
		if m.focusedArea == 1 && m.cursor == i {
			cursor = ">"
		}

		checked := "☐"
		if option.Selected {
			checked = "☑"
		}

		var optionStyle lipgloss.Style
		if m.focusedArea == 1 && m.cursor == i {
			optionStyle = styles.FocusedStyle
		} else if option.Selected {
			optionStyle = styles.CheckedStyle
		} else {
			optionStyle = styles.UnselectedStyle
		}

		choice := lipgloss.NewStyle().
			Foreground(styles.ColorTextMuted).
			Render(cursor) + " " +
			optionStyle.Render(checked+" "+option.Name)

		description := lipgloss.NewStyle().
			Foreground(styles.ColorTextMuted).
			Margin(0, 0, 0, 4).
			Render(option.Description)

		optionsList += choice + "\n" + description + "\n"
	}

	// Instructions
	instructions := lipgloss.NewStyle().
		Foreground(styles.ColorTextMuted).
		Margin(1, 0, 0, 0).
		Render("Tab: switch areas • Space: toggle features")

	return subtitle + "\n\n" + 
		   nameLabel + "\n" + nameBox + "\n\n" + 
		   optionsLabel + "\n" + optionsList + "\n" + 
		   instructions
}

func getOptionsForAppType(appType models.AppType) []AppConfigOption {
	switch appType {
	case models.AppTypeReact:
		return []AppConfigOption{
			{"google-auth", "Google OAuth", "Authentication with Google Sign-In", false},
			{"stripe", "Stripe Integration", "Payment processing with Stripe", false},
			{"tailwind", "Tailwind CSS", "Utility-first CSS framework", true},
			{"shadcn", "Shadcn/ui", "Beautiful component library", false},
			{"router", "React Router", "Client-side routing", true},
		}
	case models.AppTypeNext:
		return []AppConfigOption{
			{"auth-js", "Auth.js", "Complete authentication solution", false},
			{"stripe", "Stripe Integration", "Payment processing with Stripe", false},
			{"tailwind", "Tailwind CSS", "Utility-first CSS framework", true},
			{"app-router", "App Router", "New Next.js 13+ app directory", true},
			{"vercel", "Vercel Deployment", "Optimized for Vercel hosting", false},
		}
	case models.AppTypeTanStack:
		return []AppConfigOption{
			{"query", "TanStack Query", "Powerful data synchronization", true},
			{"router", "TanStack Router", "Type-safe router", true},
			{"tailwind", "Tailwind CSS", "Utility-first CSS framework", true},
			{"auth", "Authentication", "Built-in auth system", false},
		}
	case models.AppTypeExpo:
		return []AppConfigOption{
			{"expo-router", "Expo Router", "File-based routing for React Native", true},
			{"tamagui", "Tamagui", "Universal UI system", false},
			{"dev-tools", "Expo Dev Tools", "Enhanced development experience", true},
			{"notifications", "Push Notifications", "Expo notifications service", false},
		}
	case models.AppTypeNest:
		return []AppConfigOption{
			{"prisma", "Prisma ORM", "Type-safe database access", false},
			{"graphql", "GraphQL", "API with GraphQL", false},
			{"auth", "JWT Authentication", "JSON Web Token authentication", false},
			{"swagger", "Swagger/OpenAPI", "API documentation", true},
			{"validation", "Class Validator", "Request validation", true},
		}
	case models.AppTypeBasicNode:
		return []AppConfigOption{
			{"express", "Express.js", "Fast web framework", true},
			{"fastify", "Fastify", "Fast and low overhead web framework", false},
			{"typescript", "TypeScript", "Static type checking", true},
			{"auth", "Authentication", "Basic auth middleware", false},
		}
	default:
		return []AppConfigOption{}
	}
}

func getDefaultAppName(appType models.AppType) string {
	switch appType {
	case models.AppTypeReact:
		return "web"
	case models.AppTypeNext:
		return "web"
	case models.AppTypeTanStack:
		return "web"
	case models.AppTypeExpo:
		return "mobile"
	case models.AppTypeNest:
		return "api"
	case models.AppTypeBasicNode:
		return "api"
	default:
		return "app"
	}
}

func isValidAppNameChar(char string) bool {
	if len(char) != 1 {
		return false
	}
	c := char[0]
	return (c >= 'a' && c <= 'z') || 
		   (c >= 'A' && c <= 'Z') || 
		   (c >= '0' && c <= '9') || 
		   c == '-' || c == '_'
}

type AppConfigCompleteMsg struct {
	AppName string
	Options map[string]interface{}
}