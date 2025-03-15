package game

import (
	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/card"
	"github.com/coopersmall/subswag/domain/user"
	"github.com/coopersmall/subswag/utils"
)

const (
	RoundLimit       = 15
	RoundTimer       = 15
	StartingHandSize = 3
	HandLimit        = 5
	DeckLimit        = 52
	DrawLimit        = 1
)

type PlayerState struct {
	User                 user.UserID
	Hand                 []card.SerializableCardID // Max size of 5
	Deck                 []card.SerializableCardID // Max size of 52
	DiscardedCards       []card.SerializableCardID // Track discarded cards
	Points               int                       // Accumulated points from winning Wars
	SelectedCard         *Position                 // Current card selection for War
	RevealedCards        map[Position]bool         // Track which cards this player has seen
	EmptySpaces          []Position                // Track empty spaces on the board
	LastPlacedCard       card.SerializableCardID   // For auto-reveal if time runs out
	LastDrawnCard        card.SerializableCardID   // For auto-reveal if time runs out
	LastDiscardedCard    card.SerializableCardID   // For auto-reveal if time runs out
	HasDrawnThisTurn     bool                      // Track if player has drawn this turn
	HasSwappedThisTurn   bool                      // Track if player has swapped this turn
	HasDiscardedThisTurn bool                      // Track if player has discarded this turn
}

func NewPlayerState(
	userId user.UserID,
	deck *card.SerializableDeck,
) PlayerState {
	shuffled := card.ShuffleCards(deck.CardIDs)
	hand := shuffled[:StartingHandSize]
	return PlayerState{
		User:           userId,
		Deck:           shuffled,
		Hand:           hand,
		Points:         0,
		RevealedCards:  make(map[Position]bool),
		EmptySpaces:    make([]Position, 0),
		DiscardedCards: make([]card.SerializableCardID, 0),
	}
}

type Position struct {
	X int // 0-3
	Y int // 0-3
}

type BoardSpace struct {
	Card     card.SerializableCardID
	Revealed bool
	Owner    user.UserID // Track who placed the card
}

type GamePhase string

const (
	PhaseSetup      GamePhase = "setup"
	PhaseCardAction GamePhase = "card_action" // Player can swap/discard/draw
	PhaseReveal     GamePhase = "reveal"      // 15-second selection phase
	PhaseWar        GamePhase = "war"         // Card comparison and effect resolution
	PhaseCleanup    GamePhase = "cleanup"     // Clear revealed cards, check game end
)

type GameStateID utils.ID

func NewGameStateID() GameStateID {
	return GameStateID(utils.NewID())
}

type GameState struct {
	ID GameStateID
	GameStateData
	*domain.Metadata
}

type GameStateData struct {
	Players     [2]PlayerState
	RoundNumber int

	GamePhase

	// Board state tracking
	BoardState

	// Effect resolution tracking
	EffectsState

	// Game completion
	CompletionState

	// Rules
	Rules
}

type BoardState struct {
	Board         [4][4]BoardSpace
	ClearedSpaces map[Position]bool
}

type EffectsState struct {
	ActiveEffectsStack []EffectContext
	EffectsStack       []EffectContext
}

//	type EffectResolution struct {
//		Effect  card.CardEffect
//		Context EffectContext
//	}
type CompletionState struct {
	IsComplete bool
	Winner     *user.User
}

type Rules struct {
	RoundLimit       int
	RoundTimer       int
	HandLimit        int
	StartingHandSize int
	DeckLimit        int
	DrawLimit        int
}

type EffectContext struct {
	Trigger        Position                // Position of triggering card
	Target         Position                // Position of target card (if any)
	Source         card.SerializableCardID // Card causing the effect
	Activator      user.UserID             // Player who activated the effect
	PhaseTriggered GamePhase               // Phase when effect was triggered
}

func NewGameState(players [2]PlayerState) *GameState {
	board := [4][4]BoardSpace{}

	// Start after initial hand (0,1,2)
	p1CardIndex := 3
	p2CardIndex := 3

	// Use a boolean flag to track turns
	isPlayer1Turn := true

	// Populate board alternating between players
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			if isPlayer1Turn {
				board[y][x] = BoardSpace{
					Card:     players[0].Deck[p1CardIndex],
					Revealed: false,
					Owner:    players[0].User,
				}
				p1CardIndex++
			} else {
				board[y][x] = BoardSpace{
					Card:     players[1].Deck[p2CardIndex],
					Revealed: false,
					Owner:    players[1].User,
				}
				p2CardIndex++
			}
			// Flip the flag after each card placement
			isPlayer1Turn = !isPlayer1Turn
		}
	}

	// Update players' decks to remove placed cards
	players[0].Deck = players[0].Deck[p1CardIndex:]
	players[1].Deck = players[1].Deck[p2CardIndex:]

	return &GameState{
		ID: NewGameStateID(),
		GameStateData: GameStateData{
			Players:   players,
			GamePhase: PhaseSetup,
			BoardState: BoardState{
				Board:         board,
				ClearedSpaces: make(map[Position]bool),
			},
			EffectsState: EffectsState{
				ActiveEffectsStack: make([]EffectContext, 0),
			},
			CompletionState: CompletionState{
				IsComplete: false,
				Winner:     nil,
			},
			Rules: Rules{
				RoundLimit:       RoundLimit,
				RoundTimer:       RoundTimer,
				HandLimit:        HandLimit,
				StartingHandSize: StartingHandSize,
				DeckLimit:        DeckLimit,
				DrawLimit:        DrawLimit,
			},
		},
		Metadata: domain.NewMetadata(),
	}
}

func NewEffectContext(
	trigger Position,
	source card.SerializableCardID,
	activator user.UserID,
	phase GamePhase,
) EffectContext {
	return EffectContext{
		Trigger:        trigger,
		Source:         source,
		Activator:      activator,
		PhaseTriggered: phase,
	}
}

type GameStateVersionID utils.ID

func NewGameStateVersionID() GameStateVersionID {
	return GameStateVersionID(utils.NewID())
}

type GameStateVersion struct {
	ID       GameStateVersionID
	State    *GameState
	Metadata *domain.Metadata
}

func NewGameStateVersion(state *GameState) *GameStateVersion {
	return &GameStateVersion{
		ID:       NewGameStateVersionID(),
		State:    state,
		Metadata: domain.NewMetadata(),
	}
}
