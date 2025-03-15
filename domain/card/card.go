//ts:ignore
package card

import (
	"github.com/coopersmall/subswag/utils"
)

type CardSuite string

const (
	CardSuiteHearts   CardSuite = "hearts"
	CardSuiteDiamonds CardSuite = "diamonds"
	CardSuiteClubs    CardSuite = "clubs"
	CardSuiteSpades   CardSuite = "spades"
)

type CardRarity string

const (
	CardRarityCommon    CardRarity = "common"
	CardRarityRare      CardRarity = "rare"
	CardRarityEpic      CardRarity = "epic"
	CardRarityLegendary CardRarity = "legendary"
)

type CardTribe string

const (
	CardTribeMilitary CardTribe = "military"
	CardTribeMagic    CardTribe = "magic"
	CardTribeTech     CardTribe = "tech"
	CardTribeNature   CardTribe = "nature"
)

type CardEffectType string

const (
	CardEffectTypeDraw    CardEffectType = "draw"
	CardEffectTypeSwap    CardEffectType = "swap"
	CardEffectTypePlace   CardEffectType = "place"
	CardEffectTypeReveal  CardEffectType = "reveal"
	CardEffectTypeWar     CardEffectType = "war"
	CardEffectTypeDiscard CardEffectType = "discard"
)

type SerializableCardType string

const (
	SerializableCardTypeFace   SerializableCardType = "face"
	SerializableCardTypeNumber SerializableCardType = "number"
)

type SerializableCardID utils.ID

func NewSerializableCardID() SerializableCardID {
	return SerializableCardID(utils.NewID())
}

type SerializableDeckID utils.ID

func NewSerializableDeckID() SerializableDeckID {
	return SerializableDeckID(utils.NewID())
}
