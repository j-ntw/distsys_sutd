# Lamport’s Shared Priority Queue (Protocol)

## To enter the critical section
- Stamp the request to enter with current time T
- Add request to Qi
- Broadcast Req(T) to all other machines
- Wait until
    - Receive all replies from all other machines
    - Req(T) reaches the front of Qi

## On receipt of req(T) from some machine
- Add req(T) to Qi
- Check whether any reply is pending for an earlier request req(T’) in Qi
    - If any such reply is pending, then hold reply to req(T)

        - Note that the case signifies that the respective machine (for which the reply is not received) might not be aware of req(T’) and hence, has not replied. Thus, this particular feature of the protocol ensures property #2 (i.e. each machine is aware of all earlier requests). 
    - Otherwise, reply to req(T)
