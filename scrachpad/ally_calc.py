#!/bin/env python3
# from typing_extensions import Annotated
from enum import Enum
from pydantic import BaseModel
from typing import List, Dict, Any
import json

friendly_names = [
    "PaduriZA",
    "PolarBear",
    "CtrlAltDelicious",
    "Klaatu",
    "TheJetFlyer",
    "_PrincessMistress",
    "Captain Flint",
]
enemy_names = [
# weast_names = [
    "littler",
    "chocolo",
    "Korn",
    "Macheta",
    "ikbenjamin07",
    "BlackThorn",
    "Wizzdom",
]

enemy_names = [
    "Osric",
    "Ender",
    "SangaJanga",
    "Captain Flint",
    "Erasmas",
    "Evbogue",
]



# enum for tech:
class TechKind(Enum):
    Banking = 0
    Experimentation = 1
    Manufacturing = 2
    Range = 3
    Weapons = 5


class Tech(BaseModel):
    kind: TechKind
    level: int


class Player(BaseModel):
    home: int
    uid: int
    alias: str
    avatar: int
    race: List[int]
    color: int
    shape: int
    totalStars: int
    totalFleets: int
    totalStrength: int
    totalEconomy: int
    totalIndustry: int
    totalScience: int
    acceptedVassal: int
    offersOfFealty: List[Any]
    vassals: Dict[str, Any]
    karmaToGive: int
    ready: int
    missedTurns: int
    conceded: int
    ai: int
    regard: int
    tech: Dict[int, Tech]

    @property
    def ships_per_cycle(self):
        return self.totalIndustry * (4 + self.tech[TechKind.Manufacturing.value].level)

    @property
    def ships_per_hour(self):
        return self.ships_per_cycle / 20
    
    @property
    def next_payment(self):
        return self.totalEconomy * (10 + self.tech[TechKind.Banking.value].level * 2)

    def ships_in_cycles(self, cycles: int):
        return self.totalStrength +  self.ships_per_cycle * cycles 


class OrderPayload(BaseModel):
    players: Dict[int, Player]


class Players(BaseModel):
    players: Dict[int, Player]

def combat_simulation(def_ships:int, def_wep:int, att_ships:int, att_wep:int) -> int:
    def_wep += 1

    while def_ships > 0:
        att_ships -= def_wep
        def_ships -= att_wep

    return att_ships
    
    

# if __name__ == "__main__":

cycles = 3
friendly_manu = 9


file_path = "/home/karlis/proj/np-notify/tmp/order.json"
with open(file_path, "r") as file:
    data = json.load(file)
order = OrderPayload(**data[1])

# friendly_names += weast_names

allys: list[Player] = []
enemys: list[Player] = []
for player in order.players.values():
    if player.alias in friendly_names:
        player.tech[TechKind.Manufacturing.value].level = friendly_manu
        allys.append(player)
    elif player.alias in enemy_names:
        enemys.append(player)

def ships_in_x_cycles(player_list: list[Player]):
    return sum(
        [player.ships_in_cycles(cycles) for player in player_list]
    )

ally_total_ships = ships_in_x_cycles(allys)
enemy_total_ships = ships_in_x_cycles(enemys)



ally_wep_level = 11
enemy_wep_level = 12

ally_ships = combat_simulation(enemy_total_ships, enemy_wep_level, ally_total_ships, ally_wep_level)

print(f" assuming:ally_wep: {ally_wep_level} enemy_wep: {enemy_wep_level}, cycles: {cycles}, ally_manu: {friendly_manu}")
print(f"Ally ships in {cycles} cycle: {ally_total_ships}")
print(f"Enemy ships in {cycles} cycle: {enemy_total_ships}")


if ally_ships < 0:
    print(f"Additional ships neeeded to win: {abs(ally_ships)}")
    print(f"Total ships needed: {ally_total_ships + abs(ally_ships)}")
else:
    print(f"Ally ships after combat: {ally_ships}")

# print(f"Ally ships after combat: {ally_ships}")

pass
