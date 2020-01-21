using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class KeyboardGeneration : MonoBehaviour {

  public Transform textObject;
  public List<Transform> keyboards;

  void Start() {
    for (int i=0; i<keyboards.Count; i++) {
      Transform currentKeyboard = keyboards[i];
      for (int j=0; j<currentKeyboard.childCount; j++) {
        Transform currentKey = currentKeyboard.GetChild(j);
        GameObject inst = Instantiate(textObject.gameObject, currentKey.GetComponent<Collider>().bounds.center, currentKey.rotation, currentKey);
        inst.transform.GetChild(0).GetComponent<TextMesh>().text = currentKey.name;
      }
    }
  }

}
