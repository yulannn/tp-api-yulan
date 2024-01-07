import base64
import requests

def get_access_token(client_id, client_secret):
    client_creds = f"{client_id}:{client_secret}"
    client_creds_b64 = base64.b64encode(client_creds.encode()).decode()

    # URL pour obtenir le token d'accès
    token_url = "https://accounts.spotify.com/api/token"

    # Paramètres de la requête POST
    token_data = {
        "grant_type": "client_credentials"
    }

    # En-têtes de la requête
    token_headers = {
        "Authorization": f"Basic {client_creds_b64}"
    }

    # Envoi de la requête POST pour obtenir le token d'accès
    response = requests.post(token_url, data=token_data, headers=token_headers)

    # Vérification de la réponse
    if response.status_code == 200:
        token_response = response.json()
        access_token = token_response.get("access_token")
        return access_token
    else:
        print("Échec de l'obtention du token d'accès:", response.text)
        return None


client_id = "f27559efda8e43f0afc293f6ef47bdd6"
client_secret = "edf9ebfe1b8549e6b7de679c562ea349"

access_token = get_access_token(client_id, client_secret)
if access_token:
    print("Token d'accès obtenu avec succès:", access_token)
    # Utilisez ce token pour effectuer des requêtes vers l'API Spotify
    # Par exemple, pour obtenir des informations sur un artiste
    # headers = {"Authorization": f"Bearer {access_token}"}
    # response = requests.get("https://api.spotify.com/v1/artists/ARTIST_ID", headers=headers)
    # print(response.json())
