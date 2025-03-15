package card

import (
	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/booleanexpression"
)

type Card interface {
	GetID() SerializableCardID
	GetType() SerializableCardType
	GetMetadata() *domain.Metadata
}

type SerializableFaceCard struct {
	SerializableCardBaseData `json:",inline" validate:"required" tstype:",extends"`
	Type                     SerializableCardType `json:"type" validate:"required,eq=face" tstype:"'face'"`
	Face                     string               `json:"face" validate:"required" tstype:"string"`
}

func (f *SerializableFaceCard) GetType() SerializableCardType {
	return f.Type
}

type SerializableNumberCard struct {
	SerializableCardBaseData `json:",inline" validate:"required" tstype:",extends"`
	Type                     SerializableCardType `json:"type" validate:"required,eq=number" tstype:"'number'"`
	Number                   int                  `json:"number" validate:"required" tstype:"number"`
}

func (n *SerializableNumberCard) GetType() SerializableCardType {
	return n.Type
}

type SerializableCardBaseData struct {
	ID                            SerializableCardID `json:"id" validate:"required,gt=0"`
	SerializableCardData          `json:",inline" validate:"required" tstype:",extends"`
	SeralizableFaceCardEffectData `json:",inline" validate:"required" tstype:",extends"`
	Metadata                      *domain.Metadata `json:"metadata" validate:"required" tstype:"Metadata"`
}

func (c *SerializableCardBaseData) GetID() SerializableCardID {
	return c.ID
}

func (c *SerializableCardBaseData) GetMetadata() *domain.Metadata {
	return c.Metadata
}

type SerializableCardData struct {
	ArtworkURL string     `json:"artwork_url" validate:"required" tstype:"string"`
	Suite      CardSuite  `json:"suite" validate:"required" tstype:"Suite"`
	Rarity     CardRarity `json:"rarity" validate:"required" tstype:"Rarity"`
	Tribe      CardTribe  `json:"tribe" validate:"required" tstype:"string"`
}

type SeralizableFaceCardEffectData struct {
	OnDrawEffects    []CardEffect `json:"on_draw_effects" validate:"required" tstype:"Array<OnDrawEffectAttributes>"`
	OnRevealEffects  []CardEffect `json:"on_reveal_effects" validate:"required" tstype:"Array<OnRevealEffectAttributes>"`
	OnWarEffects     []CardEffect `json:"on_war_effects" validate:"required" tstype:"Array<OnWarEffectAttributes>"`
	OnDiscardEffects []CardEffect `json:"on_discard_effects" validate:"required" tstype:"Array<OnDiscardEffectAttributes>"`
}

type CardEffect struct {
	Condition                    booleanexpression.BooleanExpression `json:"condition" validate:"required" tstype:"string"`
	GainValueEffectAttributes    []GainValueEffectAttributes         `json:"gain_value_effects" validate:"required" tstype:"Array<GainValueEffectAttributes>"`
	GainPointsEffectAttributes   []GainPointsEffectAttributes        `json:"gain_points_effects" validate:"required" tstype:"Array<GainPointsEffectAttributes>"`
	SwapPositionEffectAttributes []SwapPositionEffectAttributes      `json:"swap_position_effects" validate:"required" tstype:"Array<SwapPositionEffectAttributes>"`
	LoseValueEffectAttributes    []LoseValueEffectAttributes         `json:"lose_value_effects" validate:"required" tstype:"Array<LoseValueEffectAttributes>"`
	LosePointsEffectAttributes   []LoseEffectAttributes              `json:"lose_effects" validate:"required" tstype:"Array<LoseEffectAttributes>"`
	DrawEffectAttributes         []DrawEffectAttributes              `json:"draw_effects" validate:"required" tstype:"Array<DrawEffectAttributes>"`
	DiscardEffectAttributes      []DiscardEffectAttributes           `json:"discard_effects" validate:"required" tstype:"Array<DiscardEffectAttributes>"`
}

type EffectAttributeBase struct {
	IsQuickEffect bool `json:"is_quick_effect" validate:"required" tstype:"boolean"`
	ForOpponent   bool `json:"for_opponent" validate:"required" tstype:"boolean"`
}

type GainValueEffectAttributes struct {
	EffectAttributeBase `json:",inline" validate:"required" tstype:",extends"`
	Amount              int `json:"amount" validate:"required" tstype:"number"`
}

type GainPointsEffectAttributes struct {
	EffectAttributeBase `json:",inline" validate:"required" tstype:",extends"`
	Amount              int `json:"amount" validate:"required" tstype:"number"`
}

type SwapPositionEffectAttributes struct {
	EffectAttributeBase `json:",inline" validate:"required" tstype:",extends"`
}

type LoseValueEffectAttributes struct {
	EffectAttributeBase `json:",inline" validate:"required" tstype:",extends"`
	Amount              int `json:"amount" validate:"required" tstype:"number"`
}

type LoseEffectAttributes struct {
	EffectAttributeBase `json:",inline" validate:"required" tstype:",extends"`
	Amount              int `json:"amount" validate:"required" tstype:"number"`
}

type DrawEffectAttributes struct {
	EffectAttributeBase `json:",inline" validate:"required" tstype:",extends"`
	Amount              int `json:"amount" validate:"required" tstype:"number"`
}

type DiscardEffectAttributes struct {
	EffectAttributeBase `json:",inline" validate:"required" tstype:",extends"`
	Amount              int `json:"amount" validate:"required" tstype:"number"`
}
