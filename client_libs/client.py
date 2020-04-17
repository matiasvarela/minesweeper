import requests
import time

class MinesweeperAPI:
    baseURL = 'http://minesweeper-api.appspot.com'
         
    def createGame(self, rows, columns, bombs):
        return requests.post(
            url=self.baseURL+'/games',
            json={
                "rows": rows,
                "columns": columns,
                "bombs": bombs
            }
        ).json()

    def getGame(self, gameId):
        return requests.get(url=self.baseURL+'/games/'+gameId).json()

    def markSquare(self, gameId, row, column):
        return requests.put(
            url=self.baseURL+'/games/'+gameId+'/mark-square',
            json={
                "row": row,
                "column": column
            }).json()

    def playSquare(self, gameId, row, column):
        return requests.put(
            url=self.baseURL+'/games/'+gameId+'/play-square',
            json={
                "row": row,
                "column": column
            }).json()