from fastapi import FastAPI, HTTPException
from pydantic import BaseModel, Field, ValidationError
import requests

app = FastAPI()

class Data(BaseModel):
    amount: float = Field(...)
    currency: str = Field(...)
    

@app.post("/convert")
def convert(data: Data):
    try:
        if data.currency == "USD":
            return data
        elif data.currency != "RUB":
            requestBody = {
                "amount": data.amount,
                "currency": data.currency
            }
            response = requests.post("http://go-service:8080/convert", json=requestBody)

            if response.status_code != 200:
                if response.status_code == 400:
                    message = "Bad Request"
                else:
                    message = "Internal Server Error"
                raise HTTPException(status_code=response.status_code, detail=message)
            
            result = response.json()

            return {"total in " + data.currency : result["total"]}
        else:
            response = requests.get("https://open.er-api.com/v6/latest/USD")

            if response.status_code != 200:
                raise HTTPException(status_code=500, detail="Internal Server Error")
            
            result = response.json()

            return {"total in RUB" : data.amount * result["rates"]["RUB"]}
    except ValidationError as e:
        raise HTTPException(status_code=422, detail=e.errors())


