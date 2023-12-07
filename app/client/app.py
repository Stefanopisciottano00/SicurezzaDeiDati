from flask import Flask, render_template, request
import subprocess
import json

string1 = 'peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/../../my-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C hospitalchannel -n Simple --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/../../my-network/organizations/peerOrganizations/hospital1.example.com/peers/peer0.hospital1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/../../my-network/organizations/peerOrganizations/hospital2.example.com/peers/peer0.hospital2.example.com/tls/ca.crt" -c "{"function":"Increment","Args":[]}"'

app = Flask(__name__)

@app.route('/')
def index():
    return render_template('index.html')

@app.route('/invoke', methods=['POST'])
def invoke():
    # Esegui l'operazione di invoke usando subprocess
    invoke_command = string1

    print("-------------------------------------VERDURE GRIGLIATE--------------------------")

    a = subprocess.run(invoke_command, shell=True)

    print(a)
    return render_template("index.html")

@app.route('/query')
def query():
    # Esegui l'operazione di query usando subprocess
    query_command = 'peer chaincode query -C hospitalchannel -n Simple -c "{"Args":["GetValue",""]}"'
    result = subprocess.run(query_command, shell=True, capture_output=True, text=True)
    return f"Risultato della query: {result.stdout}"

if __name__ == '__main__':
    app.run(debug=True)
