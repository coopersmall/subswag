package game

import (
	"context"

	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/card"
	"github.com/coopersmall/subswag/domain/game"
	"github.com/coopersmall/subswag/domain/user"
	"github.com/coopersmall/subswag/repos"
	"github.com/coopersmall/subswag/services"
)

type GameRunnerService struct {
	gameStateRepo        repos.IGameStateRepo
	gameStateVersionRepo repos.IGameStateVersionRepo
	decksService         services.IDecksService
	usersService         services.IUsersService
}

type StartGameRequest struct {
	Player1 struct {
		UserID user.UserID
		DeckID card.SerializableDeckID
	}
	Player2 struct {
		UserID user.UserID
		DeckID card.SerializableDeckID
	}
}

func (s *GameRunnerService) InitializeGame(
	ctx context.Context,
	req StartGameRequest,
) (*game.GameState, error) {
	player1, err := s.usersService.GetUser(ctx, req.Player1.UserID)
	if err != nil {
		return nil, err
	}
	player1Deck, err := s.decksService.GetDeck(ctx, req.Player1.DeckID)
	if err != nil {
		return nil, err
	}
	player2, err := s.usersService.GetUser(ctx, req.Player2.UserID)
	if err != nil {
		return nil, err
	}
	player2Deck, err := s.decksService.GetDeck(ctx, req.Player2.DeckID)
	if err != nil {
		return nil, err
	}
	player1State := game.NewPlayerState(player1.ID, player1Deck)
	player2State := game.NewPlayerState(player2.ID, player2Deck)

	gameState := game.NewGameState([2]game.PlayerState{player1State, player2State})

	if err := domain.Validate(gameState); err != nil {
		return nil, err
	}

	if err := s.gameStateRepo.Create(ctx, gameState); err != nil {
		return nil, err
	}

	if err := s.gameStateVersionRepo.Create(ctx, game.NewGameStateVersion(gameState)); err != nil {
		return nil, err
	}

	return gameState, nil
}

type IGameRunnerContext interface {
	GetGameStateData() game.GameStateData
	GetPlayerState(playerIndex int) game.PlayerState
	UpdatePlayerState(playerIndex int, playerState game.PlayerState)
	UpdateState(ctx context.Context, state game.GameStateData)
	Commit(ctx context.Context) error
}

type GameRunnerContext struct {
	gameState            *game.GameState
	gameStateRepo        repos.IGameStateRepo
	gameStateVersionRepo repos.IGameStateVersionRepo
}

func (s *GameRunnerContext) GetGameStateData() game.GameStateData {
	return s.gameState.GameStateData
}

func (s *GameRunnerContext) GetPlayerState(playerIndex int) game.PlayerState {
	return s.gameState.Players[playerIndex]
}

func (s *GameRunnerContext) GetBoardRunnerContext() game.BoardState {
	return s.gameState.BoardState
}

func (s *GameRunnerContext) UpdatePlayerState(playerIndex int, playerState game.PlayerState) {
	s.gameState.Players[playerIndex] = playerState
}

func (s *GameRunnerContext) UpdateState(ctx context.Context, state game.GameStateData) {
	s.gameState.GameStateData = state
}

func (s *GameRunnerContext) Commit(ctx context.Context) error {
	if err := s.gameStateRepo.Update(ctx, s.gameState); err != nil {
		return err
	}
	if err := s.gameStateVersionRepo.Create(ctx, game.NewGameStateVersion(s.gameState)); err != nil {
		return err
	}
	return nil
}

type IBoardRunnerContext interface {
	PlaceCard(ctx context.Context, cardId card.SerializableCardID, playerId user.UserID, position game.Position) bool
	RevealCard(ctx context.Context, playerId user.UserID, position game.Position) bool
	GetBoardState() game.BoardState
	ClearSpace(ctx context.Context, position game.Position) bool
	IsEmptySpace(position game.Position) bool
}

type BoardRunnerContext struct {
	getBoardState    func() game.BoardState
	getGamePhase     func() game.GamePhase
	updateBoardState func(ctx context.Context, board game.BoardState)
	do               func(card.SerializableCardID, card.CardEffectType, func())
}

func NewBoardRunnerContext(
	getBoardState func() game.BoardState,
	getGamePhase func() game.GamePhase,
	updateBoardState func(ctx context.Context, board game.BoardState),
	do func(card.SerializableCardID, card.CardEffectType, func()),
) IBoardRunnerContext {
	return &BoardRunnerContext{
		getBoardState:    getBoardState,
		getGamePhase:     getGamePhase,
		updateBoardState: updateBoardState,
		do:               do,
	}
}

func (s *BoardRunnerContext) GetBoardState() game.BoardState {
	return s.getBoardState()
}

func (s *BoardRunnerContext) PlaceCard(
	ctx context.Context,
	cardId card.SerializableCardID,
	playerId user.UserID,
	position game.Position,
) bool {
	onUse := func() {
		boardState := s.getBoardState()
		board := boardState.Board
		if board[position.X][position.Y].Card != 0 {
			return
		}
		for i, c := range board {
			for j, d := range c {
				if i == position.X && j == position.Y {
					if board[i][j].Card != 0 {
						return
					}
					if d.Owner != playerId {
						return
					}
					d.Card = cardId
					board[i][j] = d
					break
				}
			}
		}
		boardState.Board = board
		s.updateBoardState(ctx, boardState)
	}
	s.do(cardId, card.CardEffectTypePlace, onUse)
	return true
}

func (s *BoardRunnerContext) RevealCard(
	ctx context.Context,
	playerId user.UserID,
	position game.Position,
) bool {
	boardState := s.getBoardState()
	board := boardState.Board
	for i, c := range board {
		for j, d := range c {
			if i == position.X && j == position.Y {
				onUse := func() {
					if d.Owner != playerId {
						return
					}
					d.Revealed = true
					board[i][j] = d
					boardState.Board = board
					s.updateBoardState(ctx, boardState)
					return
				}
				s.do(d.Card, card.CardEffectTypeReveal, onUse)
				break
			}
		}
	}
	return true
}

func (s *BoardRunnerContext) ClearSpace(ctx context.Context, position game.Position) bool {
	cardId := s.getBoardState().Board[position.X][position.Y].Card
	onUse := func() {
		boardState := s.getBoardState()
		boardState.ClearedSpaces[position] = true
		boardState.Board[position.X][position.Y] = game.BoardSpace{}
		s.updateBoardState(ctx, boardState)
	}
	s.do(cardId, card.CardEffectTypeReveal, onUse)
	return true
}

func (s *BoardRunnerContext) IsEmptySpace(position game.Position) bool {
	boardState := s.getBoardState()
	return boardState.Board[position.X][position.Y].Card == 0
}

type IPlayerRunnerContext interface {
	GetPlayerID() user.UserID
	DrawCard(ctx context.Context)
	SwapCard(ctx context.Context, cardId card.SerializableCardID, position game.Position)
	DiscardCard(ctx context.Context, cardId card.SerializableCardID)
	RevealCard(ctx context.Context, position game.Position)
	AddPoints(ctx context.Context, points int)
	SubtractPoints(ctx context.Context, points int)
	ShuffleHand()
	ShuffleDeck()
}

type PlayerRunnerContext struct {
	playerId          user.UserID
	getRules          func() game.Rules
	getDeck           func() IDecKRunnerContext
	getHand           func() IHandRunnerContext
	getBoard          func() IBoardRunnerContext
	getPlayerState    func() game.PlayerState
	updatePlayerState func(playerState game.PlayerState)
	do                func(card.SerializableCardID, card.CardEffectType, func())
}

func NewPlayerRunnerContext(
	playerId user.UserID,
	getRules func() game.Rules,
	getDeck func() IDecKRunnerContext,
	getHand func() IHandRunnerContext,
	getBoard func() IBoardRunnerContext,
	getPlayerState func() game.PlayerState,
	updatePlayerState func(playerState game.PlayerState),
	do func(card.SerializableCardID, card.CardEffectType, func()),

) IPlayerRunnerContext {
	return &PlayerRunnerContext{
		playerId:          playerId,
		getDeck:           getDeck,
		getPlayerState:    getPlayerState,
		getHand:           getHand,
		getRules:          getRules,
		getBoard:          getBoard,
		updatePlayerState: updatePlayerState,
		do:                do,
	}
}

func (s *PlayerRunnerContext) GetPlayerID() user.UserID {
	return s.playerId
}

func (s *PlayerRunnerContext) DrawCard(ctx context.Context) {
	topCard, ok := s.getDeck().Peek(ctx)
	if !ok {
		return
	}
	onUse := func() {
		if s.getPlayerState().HasDrawnThisTurn {
			return
		}
		hand := s.getHand().GetHand()
		rules := s.getRules()
		if len(hand) >= rules.HandLimit {
			return
		}
		s.getDeck().Pop(ctx)
		s.getHand().AddCard(topCard)
		playerState := s.getPlayerState()
		playerState.LastDrawnCard = topCard
		playerState.HasDrawnThisTurn = true
		s.updatePlayerState(playerState)
	}
	s.do(topCard, card.CardEffectTypeDraw, onUse)
}

func (s *PlayerRunnerContext) SwapCard(
	ctx context.Context,
	cardId card.SerializableCardID,
	position game.Position,
) {
	onUse := func() {
		hand := s.getHand()
		hand.RemoveCard(cardId)
		s.getBoard().PlaceCard(ctx, cardId, s.playerId, position)
		playerState := s.getPlayerState()
		playerState.LastPlacedCard = cardId
		playerState.HasSwappedThisTurn = true
		s.updatePlayerState(playerState)
	}
	s.do(cardId, card.CardEffectTypeSwap, onUse)
}

func (s *PlayerRunnerContext) AddCardFromBoardToHand(
	ctx context.Context,
	position game.Position,
) {
	boardState := s.getBoard().GetBoardState()
	cardId := boardState.Board[position.X][position.Y].Card
	if cardId == 0 {
		return
	}
	onUse := func() {
		s.getBoard().ClearSpace(ctx, position)
		s.getHand().AddCard(cardId)
	}
	s.do(cardId, card.CardEffectTypeDraw, onUse)
}

func (s *PlayerRunnerContext) DiscardCard(ctx context.Context, cardId card.SerializableCardID) {
	onUse := func() {
		s.getHand().RemoveCard(cardId)
		playerState := s.getPlayerState()
		newDiscarded := make([]card.SerializableCardID, len(playerState.DiscardedCards)+1)
		copy(newDiscarded, playerState.DiscardedCards)
		newDiscarded[len(playerState.DiscardedCards)] = cardId
		playerState.DiscardedCards = newDiscarded
		playerState.LastDiscardedCard = cardId
		playerState.HasDiscardedThisTurn = true
		s.updatePlayerState(playerState)
	}
	s.do(cardId, card.CardEffectTypeDiscard, onUse)
}

func (s *PlayerRunnerContext) AddDiscardedCardToHand(ctx context.Context, cardId card.SerializableCardID) {
	playerState := s.getPlayerState()
	newDiscarded := make([]card.SerializableCardID, 0, len(playerState.DiscardedCards)-1)
	for _, c := range playerState.DiscardedCards {
		if c != cardId {
			newDiscarded = append(newDiscarded, c)
		}
	}
	playerState.DiscardedCards = newDiscarded
	s.updatePlayerState(playerState)
	s.getHand().AddCard(cardId)
}

func (s *PlayerRunnerContext) AddDiscardedCardsToDeck(ctx context.Context) {
	playerState := s.getPlayerState()
	discarded := playerState.DiscardedCards
	playerState.DiscardedCards = nil
	s.updatePlayerState(playerState)
	s.getDeck().PushMany(ctx, discarded...)
}

func (s *PlayerRunnerContext) AddDiscardedCardToDeck(ctx context.Context, cardId card.SerializableCardID) {
	playerState := s.getPlayerState()
	newDiscarded := make([]card.SerializableCardID, 0, len(playerState.DiscardedCards)-1)
	for _, c := range playerState.DiscardedCards {
		if c != cardId {
			newDiscarded = append(newDiscarded, c)
		}
	}
	playerState.DiscardedCards = newDiscarded
	s.updatePlayerState(playerState)
	s.getDeck().Push(ctx, cardId)
}

func (s *PlayerRunnerContext) DiscardHand(ctx context.Context) {
	playerState := s.getPlayerState()
	hand := s.getHand().GetHand()
	newDiscarded := make([]card.SerializableCardID, len(playerState.DiscardedCards)+len(hand))
	copy(newDiscarded, playerState.DiscardedCards)
	copy(newDiscarded[len(playerState.DiscardedCards):], hand)
	playerState.DiscardedCards = newDiscarded
	s.updatePlayerState(playerState)
	s.getDeck().PushMany(ctx, hand...)
}

func (s *PlayerRunnerContext) AddPoints(ctx context.Context, points int) {
	playerState := s.getPlayerState()
	playerState.Points += points
	s.updatePlayerState(playerState)
}

func (s *PlayerRunnerContext) SubtractPoints(ctx context.Context, points int) {
	playerState := s.getPlayerState()
	playerState.Points -= points
	s.updatePlayerState(playerState)
}

func (s *PlayerRunnerContext) RevealCard(ctx context.Context, position game.Position) {
	s.getBoard().RevealCard(ctx, s.playerId, position)
}

func (s *PlayerRunnerContext) ShuffleHand() {
	s.getHand().ShuffleHand()
}

func (s *PlayerRunnerContext) ShuffleDeck() {
	s.getDeck().Shuffle()
}

type IDecKRunnerContext interface {
	Pop(ctx context.Context) (card.SerializableCardID, bool)
	PopMany(ctx context.Context, count int) ([]card.SerializableCardID, bool)
	PopBottom(ctx context.Context) (card.SerializableCardID, bool)
	Peek(ctx context.Context) (card.SerializableCardID, bool)
	PeekMany(ctx context.Context, count int) ([]card.SerializableCardID, bool)
	PeekBottom(ctx context.Context) (card.SerializableCardID, bool)
	Push(ctx context.Context, cardId card.SerializableCardID)
	PushMany(ctx context.Context, cardIds ...card.SerializableCardID)
	PushBottom(ctx context.Context, cardId card.SerializableCardID)
	Shuffle()
}

type DeckRunnerContext struct {
	playerIdx  int
	getDeck    func() []card.SerializableCardID
	updateDeck func(deck []card.SerializableCardID)
}

func NewDeckRunnerContext(
	playerIdx int,
	getDeck func() []card.SerializableCardID,
	updateDeck func(deck []card.SerializableCardID),
) IDecKRunnerContext {
	return &DeckRunnerContext{
		playerIdx:  playerIdx,
		getDeck:    getDeck,
		updateDeck: updateDeck,
	}
}

func (s *DeckRunnerContext) GetDeck() []card.SerializableCardID {
	return s.getDeck()
}

func (s *DeckRunnerContext) Pop(ctx context.Context) (card.SerializableCardID, bool) {
	deck := s.getDeck()
	if len(deck) == 0 {
		return 0, false
	}
	topCard := deck[len(deck)-1]
	newDeck := make([]card.SerializableCardID, len(deck)-1)
	copy(newDeck, deck[:len(deck)-1])
	s.updateDeck(newDeck)
	return topCard, true
}

func (s *DeckRunnerContext) PopMany(ctx context.Context, count int) ([]card.SerializableCardID, bool) {
	deck := s.getDeck()
	if len(deck) < count {
		return nil, false
	}
	topCards := deck[len(deck)-count:]
	newDeck := make([]card.SerializableCardID, len(deck)-count)
	copy(newDeck, deck[:len(deck)-count])
	s.updateDeck(newDeck)
	return topCards, true
}

func (s *DeckRunnerContext) PopBottom(ctx context.Context) (card.SerializableCardID, bool) {
	deck := s.getDeck()
	if len(deck) == 0 {
		return 0, false
	}
	bottomCard := deck[0]
	newDeck := make([]card.SerializableCardID, len(deck)-1)
	copy(newDeck, deck[1:])
	s.updateDeck(newDeck)
	return bottomCard, true
}

func (s *DeckRunnerContext) Peek(ctx context.Context) (card.SerializableCardID, bool) {
	deck := s.getDeck()
	if len(deck) == 0 {
		return 0, false
	}
	return deck[len(deck)-1], true
}

func (s *DeckRunnerContext) PeekMany(ctx context.Context, count int) ([]card.SerializableCardID, bool) {
	deck := s.getDeck()
	if len(deck) < count {
		return nil, false
	}
	return deck[len(deck)-count:], true
}

func (s *DeckRunnerContext) PeekBottom(ctx context.Context) (card.SerializableCardID, bool) {
	deck := s.getDeck()
	if len(deck) == 0 {
		return 0, false
	}
	return deck[0], true
}

func (s *DeckRunnerContext) Push(ctx context.Context, cardId card.SerializableCardID) {
	deck := s.getDeck()
	newDeck := make([]card.SerializableCardID, len(deck)+1)
	copy(newDeck, deck)
	newDeck[len(deck)] = cardId
	s.updateDeck(newDeck)
}

func (s *DeckRunnerContext) PushMany(ctx context.Context, cardIds ...card.SerializableCardID) {
	deck := s.getDeck()
	newDeck := make([]card.SerializableCardID, len(deck)+len(cardIds))
	copy(newDeck, deck)
	copy(newDeck[len(deck):], cardIds)
	s.updateDeck(newDeck)
}

func (s *DeckRunnerContext) PushBottom(ctx context.Context, cardId card.SerializableCardID) {
	deck := s.getDeck()
	newDeck := make([]card.SerializableCardID, len(deck)+1)
	newDeck[0] = cardId
	copy(newDeck[1:], deck)
	s.updateDeck(newDeck)
}

func (s *DeckRunnerContext) Shuffle() {
	deck := s.getDeck()
	shuffled := card.ShuffleCards(deck)
	s.updateDeck(shuffled)
}

type IHandRunnerContext interface {
	GetHand() []card.SerializableCardID
	AddCard(cardId card.SerializableCardID)
	RemoveCard(cardId card.SerializableCardID)
	ShuffleHand()
}

type HandRunnerContext struct {
	playerIdx  int
	getRules   func() game.Rules
	getHand    func() []card.SerializableCardID
	updateHand func(hand []card.SerializableCardID)
}

func NewHandRunnerContext(
	playerIdx int,
	getRules func() game.Rules,
	getHand func() []card.SerializableCardID,
	updateHand func(hand []card.SerializableCardID),
) IHandRunnerContext {
	return &HandRunnerContext{
		playerIdx:  playerIdx,
		getRules:   getRules,
		getHand:    getHand,
		updateHand: updateHand,
	}
}

func (s *HandRunnerContext) GetHand() []card.SerializableCardID {
	return s.getHand()
}

func (s *HandRunnerContext) AddCard(cardId card.SerializableCardID) {
	hand := s.getHand()
	rules := s.getRules()
	if len(hand) >= rules.HandLimit {
		return
	}
	newHand := make([]card.SerializableCardID, len(hand)+1)
	copy(newHand, hand)
	newHand[len(hand)] = cardId
	s.updateHand(newHand)
}

func (s *HandRunnerContext) RemoveCard(cardId card.SerializableCardID) {
	hand := s.getHand()
	if len(hand) == 0 {
	}
	newHand := make([]card.SerializableCardID, 0, len(hand)-1)
	for _, c := range hand {
		if c != cardId {
			newHand = append(newHand, c)
		}
	}
	s.updateHand(newHand)
}

func (s *HandRunnerContext) ShuffleHand() {
	hand := s.getHand()
	shuffled := card.ShuffleCards(hand)
	s.updateHand(shuffled)
}
